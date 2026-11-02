// Package search provides Elasticsearch-backed product search.
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

const (
	indexName               = "mall-products"
	indexRequestTimeout     = 2 * time.Second
	indexBatchSize          = model.MaxPageSize
	productNameField        = "name"
	productDescriptionField = "description"
	productCategoryField    = "categoryId"
	productShopField        = "shopId"
	textAnalyzerName        = "mall_text"
)

type Elasticsearch struct {
	endpoint string
	client   *http.Client
}

func productDocument(product model.Product) map[string]any {
	return map[string]any{
		productNameField:        product.Name,
		productDescriptionField: product.Description,
		productCategoryField:    product.CategoryID,
		productShopField:        product.ShopID,
	}
}

func New(endpoint string) *Elasticsearch {
	return &Elasticsearch{endpoint: strings.TrimRight(endpoint, "/"), client: &http.Client{Timeout: indexRequestTimeout}}
}

func (e *Elasticsearch) ensureReady(ctx context.Context) error {
	resp, err := e.request(ctx, http.MethodHead, "/"+indexName, nil)
	if err != nil {
		return err
	}
	status := resp.StatusCode
	resp.Body.Close()
	if status == http.StatusOK {
		return nil
	}
	if status != http.StatusNotFound {
		return fmt.Errorf("elasticsearch index check status %d", status)
	}
	return e.EnsureIndex(ctx)
}

// EnsureIndex creates the index with a Chinese analyzer when it does not exist.
func (e *Elasticsearch) EnsureIndex(ctx context.Context) error {
	resp, err := e.request(ctx, http.MethodPut, "/"+indexName, map[string]any{
		"settings": map[string]any{"analysis": map[string]any{"analyzer": map[string]any{textAnalyzerName: map[string]any{"type": "custom", "tokenizer": "standard", "filter": []string{"lowercase", "cjk_width", "cjk_bigram"}}}}},
		"mappings": map[string]any{"properties": map[string]any{
			productNameField: map[string]any{"type": "text", "analyzer": textAnalyzerName}, productDescriptionField: map[string]any{"type": "text", "analyzer": textAnalyzerName},
			productCategoryField: map[string]any{"type": "keyword"}, productShopField: map[string]any{"type": "keyword"},
		}},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("elasticsearch index setup status %d", resp.StatusCode)
	}
	return nil
}

func (e *Elasticsearch) request(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, e.endpoint+path, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return e.client.Do(req)
}

func (e *Elasticsearch) Index(ctx context.Context, p model.Product) error {
	if err := e.ensureReady(ctx); err != nil {
		return err
	}
	body := productDocument(p)
	resp, err := e.request(ctx, http.MethodPut, "/"+indexName+"/_doc/"+url.PathEscape(p.ID)+"?refresh=wait_for", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch index status %d", resp.StatusCode)
	}
	return nil
}

func (e *Elasticsearch) Delete(ctx context.Context, id string) error {
	resp, err := e.request(ctx, http.MethodDelete, "/"+indexName+"/_doc/"+url.PathEscape(id)+"?refresh=wait_for", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("elasticsearch delete status %d", resp.StatusCode)
	}
	return nil
}

// Reindex replaces catalog documents using a caller supplied, paginated source.
func (e *Elasticsearch) Reindex(ctx context.Context, load func(context.Context, model.ListOptions) (model.Page[model.Product], error)) error {
	if err := e.EnsureIndex(ctx); err != nil {
		return err
	}
	resp, err := e.request(ctx, http.MethodPost, "/"+indexName+"/_delete_by_query?refresh=true", map[string]any{"query": map[string]any{"match_all": map[string]any{}}})
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		resp.Body.Close()
		return fmt.Errorf("elasticsearch clear index status %d", resp.StatusCode)
	}
	resp.Body.Close()
	for pageNumber := 1; ; pageNumber++ {
		page, err := load(ctx, model.ListOptions{Page: pageNumber, PageSize: indexBatchSize})
		if err != nil {
			return err
		}
		for _, product := range page.Items {
			body := productDocument(product)
			resp, err = e.request(ctx, http.MethodPut, "/"+indexName+"/_doc/"+url.PathEscape(product.ID), body)
			if err != nil {
				return err
			}
			if resp.StatusCode >= 300 {
				resp.Body.Close()
				return fmt.Errorf("elasticsearch reindex status %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
		if int64(pageNumber*indexBatchSize) >= page.Total || len(page.Items) == 0 {
			break
		}
	}
	resp, err = e.request(ctx, http.MethodPost, "/"+indexName+"/_refresh", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch refresh status %d", resp.StatusCode)
	}
	return nil
}

func (e *Elasticsearch) Search(ctx context.Context, keyword string, opts model.ListOptions) (model.Page[model.Product], error) {
	if err := e.ensureReady(ctx); err != nil {
		return model.Page[model.Product]{}, err
	}
	query := map[string]any{"from": (opts.Page - 1) * opts.PageSize, "size": opts.PageSize, "query": map[string]any{"multi_match": map[string]any{"query": keyword, "fields": []string{productNameField + "^3", productDescriptionField}, "operator": "and", "type": "best_fields"}}}
	resp, err := e.request(ctx, http.MethodPost, "/"+indexName+"/_search", query)
	if err != nil {
		return model.Page[model.Product]{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return model.Page[model.Product]{}, fmt.Errorf("elasticsearch search status %d", resp.StatusCode)
	}
	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string `json:"_id"`
				Source struct {
					CategoryID string `json:"categoryId"`
					ShopID     string `json:"shopId"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.Page[model.Product]{}, err
	}
	page := model.Page[model.Product]{Total: result.Hits.Total.Value, Items: make([]model.Product, 0, len(result.Hits.Hits))}
	for _, hit := range result.Hits.Hits {
		page.Items = append(page.Items, model.Product{Base: model.Base{ID: hit.ID}, CategoryID: hit.Source.CategoryID, ShopID: hit.Source.ShopID})
	}
	return page, nil
}

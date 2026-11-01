package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func TestSearchUsesProductFieldsAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path != "/mall-products/_search" {
			t.Errorf("path = %s", r.URL.Path)
		}
		var query map[string]any
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			t.Fatal(err)
		}
		if query["from"] != float64(10) || query["size"] != float64(5) {
			t.Errorf("pagination = %#v", query)
		}
		_, _ = w.Write([]byte(`{"hits":{"total":{"value":1},"hits":[{"_id":"p1","_source":{"categoryId":"c1","shopId":"s1"}}]}}`))
	}))
	defer server.Close()
	page, err := New(server.URL).Search(context.Background(), "手冲", model.ListOptions{Page: 3, PageSize: 5})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != "p1" {
		t.Fatalf("page = %+v", page)
	}
}

func TestSearchEnsuresMissingIndex(t *testing.T) {
	created := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodHead:
			if !created {
				w.WriteHeader(http.StatusNotFound)
			}
		case r.Method == http.MethodPut && r.URL.Path == "/mall-products":
			created = true
		case r.Method == http.MethodPost:
			_, _ = w.Write([]byte(`{"hits":{"total":{"value":0},"hits":[]}}`))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	if _, err := New(server.URL).Search(context.Background(), "杯子", model.ListOptions{Page: 1, PageSize: 10}); err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Fatal("expected index creation")
	}
}

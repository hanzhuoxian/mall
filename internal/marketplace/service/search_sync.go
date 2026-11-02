package service

import (
	"context"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

const (
	searchSyncInterval  = time.Second
	searchSyncBatchSize = model.MaxPageSize
)

// RunSearchSync continuously delivers durable outbox events to Elasticsearch.
func (s *Service) RunSearchSync(ctx context.Context) {
	if s.search == nil {
		return
	}
	ticker := time.NewTicker(searchSyncInterval)
	defer ticker.Stop()
	for {
		if err := s.syncSearchBatch(ctx); err != nil && ctx.Err() == nil { /* next tick retries */
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Service) syncSearchBatch(ctx context.Context) error {
	events, err := s.store.PendingSearchOutbox(ctx, searchSyncBatchSize)
	if err != nil {
		return err
	}
	for _, event := range events {
		var syncErr error
		product, productErr := s.store.Product(ctx, event.ProductID, false)
		if productErr != nil {
			syncErr = s.search.Delete(ctx, event.ProductID)
		} else {
			shop, shopErr := s.store.Shop(ctx, product.ShopID, false)
			if shopErr != nil {
				syncErr = shopErr
			} else if product.Status == model.ProductStatusApproved && shop.Status == model.ShopStatusApproved {
				syncErr = s.search.Index(ctx, product)
			} else {
				syncErr = s.search.Delete(ctx, product.ID)
			}
		}
		if syncErr != nil {
			_ = s.store.IncrementSearchOutboxAttempts(ctx, event.ID)
			continue
		}
		if err = s.store.DeleteSearchOutbox(ctx, event.ID); err != nil {
			return err
		}
	}
	return nil
}

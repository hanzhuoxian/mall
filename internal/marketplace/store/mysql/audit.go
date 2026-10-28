package mysql

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (s *datastore) CreateAudit(ctx context.Context, audit *model.Audit) error {
	return normalizeError(s.db.WithContext(ctx).Create(audit).Error)
}

func (s *datastore) ListAudits(ctx context.Context, opts model.ListOptions) (model.Page[model.Audit], error) {
	return list[model.Audit](s.db.WithContext(ctx).Model(&model.Audit{}), opts)
}

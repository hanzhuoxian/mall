package controller

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	marketcontroller "github.com/hanzhuoxian/mall/internal/marketplace/controller"
	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/service"
	"github.com/hanzhuoxian/mall/internal/marketplace/store/mysql"
	"github.com/hanzhuoxian/mall/internal/pkg/admin"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// NewMarketplaceService 根据配置创建商城存储及业务服务。
func NewMarketplaceService(cfg *config.Config) (*service.Service, error) {
	db, err := cfg.MySQLOptions.NewClient()
	if err != nil {
		return nil, err
	}
	return service.NewService(mysql.NewStore(db), cfg.SandboxPayment), nil
}

// NewMarketplaceController 为商城控制器注入用户身份查询能力。
func NewMarketplaceController(
	svc *service.Service,
	users *grpcclient.UserClient,
) *marketcontroller.MarketplaceController {
	identity := func(ctx context.Context, id string) (model.Actor, error) {
		resp, err := users.GetUser(ctx, &userv1.GetUserRequest{InstanceId: id})
		if err != nil {
			return model.Actor{}, err
		}
		user := resp.GetUser()
		if user == nil || user.Status != 1 {
			return model.Actor{}, service.ErrForbidden
		}
		return model.Actor{ID: id, Admin: admin.IsAdmin(user.Username)}, nil
	}
	return marketcontroller.NewMarketplaceController(svc, identity, middleware.UserIdentifier)
}

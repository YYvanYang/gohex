package query

import (
	"context"
	"github.com/your-org/your-project/internal/application/dto"
)

// GetUserRolesQuery 获取用户角色查询
type GetUserRolesQuery struct {
	UserID string
}

type GetUserRolesHandler struct {
	userRepo port.UserRepository
	logger   Logger
	metrics  MetricsReporter
}

func (h *GetUserRolesHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*GetUserRolesQuery)

	// 1. 获取用户
	user, err := h.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// 2. 转换为 DTO
	roles := make([]string, len(user.Roles()))
	for i, role := range user.Roles() {
		roles[i] = role.String()
	}

	return &dto.UserRolesDTO{
		UserID: user.ID(),
		Roles:  roles,
	}, nil
} 
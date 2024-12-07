package query

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/your-project/internal/application/dto"
)

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	ID string
}

type GetUserHandler struct {
	userRepo port.UserRepository
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func (h *GetUserHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*GetUserQuery)

	// 1. 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:%s", query.ID)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		h.metrics.IncrementCounter("user_cache_hit")
		return cached.(*dto.UserDTO), nil
	}

	// 2. 从仓储获取
	user, err := h.userRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	// 3. 转换为 DTO
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().String(),
		Name:      user.Profile().Name(),
		Bio:       user.Profile().Bio(),
		Avatar:    user.Profile().Avatar(),
		Status:    user.Status().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	// 4. 更新缓存
	if err := h.cache.Set(ctx, cacheKey, userDTO, time.Hour); err != nil {
		h.logger.Error("failed to cache user", "error", err)
	}

	return userDTO, nil
}

// ListUsersQuery 用户列表查询
type ListUsersQuery struct {
	Page     int
	PageSize int
	Status   string
	SortBy   string
	SortDir  string
}

func (q ListUsersQuery) Offset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	return (q.Page - 1) * q.Limit()
}

func (q ListUsersQuery) Limit() int {
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	return q.PageSize
}

type ListUsersHandler struct {
	userRepo port.UserRepository
	logger   Logger
	metrics  MetricsReporter
}

func (h *ListUsersHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*ListUsersQuery)

	// 1. 获取总数
	total, err := h.userRepo.Count(ctx, query.Status)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户列表
	users, err := h.userRepo.FindAll(ctx, query.Status, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	// 3. 转换为 DTO
	items := make([]dto.UserDTO, len(users))
	for i, user := range users {
		items[i] = dto.UserDTO{
			ID:        user.ID(),
			Email:     user.Email().String(),
			Name:      user.Profile().Name(),
			Bio:       user.Profile().Bio(),
			Avatar:    user.Profile().Avatar(),
			Status:    user.Status().String(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		}
	}

	return &dto.UserListDTO{
		Total: total,
		Items: items,
	}, nil
}

type FindUserByEmailQuery struct {
	Email string
}

type FindUserByIDQuery struct {
	ID string
}

type GetUserRolesQuery struct {
	UserID string
}

type GetUserPermissionsQuery struct {
	UserID string
} 
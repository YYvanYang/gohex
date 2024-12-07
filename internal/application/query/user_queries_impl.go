package query

import (
	"context"
	"fmt"
	"time"
	"github.com/your-org/your-project/internal/application/dto"
	"github.com/your-org/your-project/internal/domain/vo"
	"github.com/your-org/your-project/pkg/errors"
)

// GetUserByIDQuery 实现 Cacheable 接口
type GetUserByIDQuery struct {
	ID string `validate:"required"`
}

func (q GetUserByIDQuery) CacheKey() string {
	return fmt.Sprintf("user:id:%s", q.ID)
}

func (q GetUserByIDQuery) TTL() time.Duration {
	return time.Hour * 24
}

type GetUserByIDHandler struct {
	userRepo port.UserRepository
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func (h *GetUserByIDHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*GetUserByIDQuery)
	
	// 1. 尝试从缓存获取
	cacheKey := query.CacheKey()
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		h.metrics.IncrementCounter("cache_hit", "type", "user")
		return cached.(*dto.UserDTO), nil
	}
	h.metrics.IncrementCounter("cache_miss", "type", "user")

	// 2. 从数据库获取
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
		Roles:     user.RoleStrings(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	// 4. 更新缓存
	if err := h.cache.Set(ctx, cacheKey, userDTO, query.TTL()); err != nil {
		h.logger.Error("failed to cache user", "error", err)
	}

	return userDTO, nil
}

// ListUsersQuery 实现 Cacheable 接口
type ListUsersQuery struct {
	Page     int    `validate:"min=1"`
	PageSize int    `validate:"min=1,max=100"`
	Status   string `validate:"omitempty,oneof=active inactive suspended deleted"`
	SortBy   string `validate:"omitempty,oneof=created_at updated_at name email"`
	SortDir  string `validate:"omitempty,oneof=asc desc"`
}

func (q ListUsersQuery) CacheKey() string {
	return fmt.Sprintf("users:page:%d:page_size:%d:status:%s:sort_by:%s:sort_dir:%s", q.Page, q.PageSize, q.Status, q.SortBy, q.SortDir)
}

func (q ListUsersQuery) TTL() time.Duration {
	return time.Hour * 24
}

func (q ListUsersQuery) Validate() error {
	if q.Page <= 0 {
		return errors.NewValidationError("page must be greater than 0")
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		return errors.NewValidationError("page size must be between 1 and 100")
	}
	if q.Status != "" && !vo.UserStatus(q.Status).IsValid() {
		return errors.NewValidationError("invalid status")
	}
	if q.SortBy != "" && !isValidSortField(q.SortBy) {
		return errors.NewValidationError("invalid sort field")
	}
	if q.SortDir != "" && !isValidSortDirection(q.SortDir) {
		return errors.NewValidationError("invalid sort direction")
	}
	return nil
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":      true,
		"email":     true,
	}
	return validFields[field]
}

func isValidSortDirection(dir string) bool {
	return dir == "asc" || dir == "desc"
}

type ListUsersHandler struct {
	userRepo port.UserRepository
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func NewListUsersHandler(
	userRepo port.UserRepository,
	cache port.Cache,
	logger Logger,
	metrics MetricsReporter,
) *ListUsersHandler {
	return &ListUsersHandler{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
		metrics:  metrics,
	}
}

func (h *ListUsersHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*ListUsersQuery)
	
	// 添加缓存键前缀，便于批量清除
	cacheKey := fmt.Sprintf("users:list:%s", query.CacheKey())
	// 1. 尝试从缓存获取
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		h.metrics.IncrementCounter("cache_hit", "type", "users")
		return cached.(*dto.UserDTO), nil
	}
	h.metrics.IncrementCounter("cache_miss", "type", "users")

	// 2. 从数据库获取
	users, err := h.userRepo.FindAll(ctx, query.Page, query.PageSize, query.Status, query.SortBy, query.SortDir)
	if err != nil {
		return nil, err
	}

	// 3. 转换为 DTO
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = &dto.UserDTO{
			ID:        user.ID(),
			Email:     user.Email().String(),
			Name:      user.Profile().Name(),
			Bio:       user.Profile().Bio(),
			Avatar:    user.Profile().Avatar(),
			Status:    user.Status().String(),
			Roles:     user.RoleStrings(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		}
	}

	// 4. 更新缓存
	if err := h.cache.Set(ctx, cacheKey, userDTOs, query.TTL()); err != nil {
		h.logger.Error("failed to cache users", "error", err)
	}

	return userDTOs, nil
}

// 添加缓存清理方法
func (h *GetUserByIDHandler) clearCache(ctx context.Context, userID string) {
	cacheKeys := []string{
		fmt.Sprintf("user:id:%s", userID),
		"users:list", // 清除列表缓存
	}
	
	for _, key := range cacheKeys {
		if err := h.cache.Delete(ctx, key); err != nil {
			h.logger.Error("failed to clear cache", "key", key, "error", err)
		}
	}
} 
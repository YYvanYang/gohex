package query

import (
	"context"
	"github.com/your-org/your-project/internal/application/dto"
)

// ValidateTokenQuery 验证令牌查询
type ValidateTokenQuery struct {
	Token string `validate:"required"`
}

type ValidateTokenHandler struct {
	tokenSvc port.TokenService
	userRepo port.UserRepository
	logger   Logger
	metrics  MetricsReporter
}

func (h *ValidateTokenHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*ValidateTokenQuery)

	// 1. 验证令牌
	claims, err := h.tokenSvc.ValidateToken(ctx, query.Token)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户
	user, err := h.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// 3. 检查用户状态
	if !user.IsActive() {
		return nil, errors.ErrInactiveUser
	}

	// 4. 返回认证信息
	return &dto.AuthInfoDTO{
		UserID:    user.ID(),
		Email:     user.Email().String(),
		Roles:     user.RoleStrings(),
		ExpiresAt: claims.ExpiresAt,
	}, nil
} 
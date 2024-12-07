package query

import (
	"context"
	"github.com/your-org/your-project/internal/application/dto"
)

// GetUserByEmailQuery 通过邮箱查询用户
type GetUserByEmailQuery struct {
	Email string `validate:"required,email"`
}

type GetUserByEmailHandler struct {
	userRepo port.UserRepository
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func (h *GetUserByEmailHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	query := q.(*GetUserByEmailQuery)

	email, err := vo.NewEmail(query.Email)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().String(),
		Name:      user.Profile().Name(),
		Bio:       user.Profile().Bio(),
		Avatar:    user.Profile().Avatar(),
		Status:    user.Status().String(),
		Roles:     user.RoleStrings(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
} 
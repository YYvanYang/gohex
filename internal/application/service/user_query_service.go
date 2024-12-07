package service

import (
	"context"
	"github.com/your-org/your-project/internal/application/dto"
	"github.com/your-org/your-project/internal/domain/vo"
)

type UserQueryService struct {
	userRepo port.UserRepository
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func NewUserQueryService(
	userRepo port.UserRepository,
	cache port.Cache,
	logger Logger,
	metrics MetricsReporter,
) *UserQueryService {
	return &UserQueryService{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
		metrics:  metrics,
	}
}

func (s *UserQueryService) FindByEmail(ctx context.Context, email string) (*dto.UserDTO, error) {
	emailVO, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		return nil, err
	}

	return s.toDTO(user), nil
}

func (s *UserQueryService) ListUsers(ctx context.Context, query ListUsersQuery) (*dto.UserListDTO, error) {
	total, err := s.userRepo.Count(ctx, query.Status)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepo.FindAll(ctx, query.Status, query.Offset(), query.Limit())
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserDTO, len(users))
	for i, user := range users {
		items[i] = *s.toDTO(user)
	}

	return &dto.UserListDTO{
		Total: total,
		Items: items,
	}, nil
}

func (s *UserQueryService) toDTO(user *aggregate.User) *dto.UserDTO {
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
	}
} 
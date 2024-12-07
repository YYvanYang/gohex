package service

import (
	"context"
	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
)

type UserService struct {
	userRepo port.UserRepository
	logger   Logger
}

func NewUserService(userRepo port.UserRepository, logger Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) ValidateUniqueEmail(ctx context.Context, email vo.Email) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrEmailAlreadyExists
	}
	return nil
}

func (s *UserService) ValidateUserStatus(ctx context.Context, userID string, expectedStatus vo.UserStatus) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status() != expectedStatus {
		return errors.ErrInvalidUserStatus
	}

	return nil
}

func (s *UserService) ValidateUserPermission(ctx context.Context, userID string, permission string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, role := range user.Roles() {
		if role.HasPermission(permission) {
			return nil
		}
	}

	return errors.ErrInsufficientPermissions
} 
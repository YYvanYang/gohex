package command

import (
	"context"
	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
)

// RegisterUserCommand 注册用户命令
type RegisterUserCommand struct {
	Email    string
	Password string
	Name     string
	Bio      string
}

type RegisterUserHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	eventBus   port.EventBus
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	registerCmd := cmd.(*RegisterUserCommand)
	
	var result struct {
		ID string `json:"id"`
	}

	err := h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 创建值对象
		email, err := vo.NewEmail(registerCmd.Email)
		if err != nil {
			return err
		}

		password, err := vo.NewPassword(registerCmd.Password)
		if err != nil {
			return err
		}

		profile, err := vo.NewUserProfile(registerCmd.Name, registerCmd.Bio)
		if err != nil {
			return err
		}

		// 2. 检查邮箱是否已存在
		exists, err := h.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			return errors.ErrEmailAlreadyExists
		}

		// 3. 创建用户聚合根
		user, err := aggregate.NewUser(email, password, profile)
		if err != nil {
			return err
		}

		// 4. 保存用户
		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		// 5. 保存事件
		if err := h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), 0); err != nil {
			return err
		}

		result.ID = user.ID()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUserProfileCommand 更新用户资料命令
type UpdateUserProfileCommand struct {
	UserID   string `validate:"required"`
	Name     string `validate:"required"`
	Bio      string `validate:"max=500"`
	Avatar   string `validate:"omitempty,url"`
	Location string
	Website  string `validate:"omitempty,url"`
}

type UpdateUserProfileHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	cache      port.Cache
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *UpdateUserProfileHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	updateCmd := cmd.(*UpdateUserProfileCommand)

	err := h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 获取用户
		user, err := h.userRepo.FindByID(ctx, updateCmd.UserID)
		if err != nil {
			return err
		}

		// 2. 创建新的资料值对象
		profile, err := vo.NewUserProfile(updateCmd.Name, updateCmd.Bio)
		if err != nil {
			return err
		}

		// 3. 更新头像
		if updateCmd.Avatar != "" {
			profile, err = profile.WithAvatar(updateCmd.Avatar)
			if err != nil {
				return err
			}
		}

		// 4. 更新位置
		profile = profile.WithLocation(updateCmd.Location)

		// 5. 更新网站
		if updateCmd.Website != "" {
			profile, err = profile.WithWebsite(updateCmd.Website)
			if err != nil {
				return err
			}
		}

		// 6. 更新用户资料
		if err := user.UpdateProfile(profile); err != nil {
			return err
		}

		// 7. 保存用户
		if err := h.userRepo.Update(ctx, user); err != nil {
			return err
		}

		// 8. 保存事件
		if err := h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), user.Version()); err != nil {
			return err
		}

		// 9. 清除缓存
		cacheKey := fmt.Sprintf("user:%s", user.ID())
		if err := h.cache.Delete(ctx, cacheKey); err != nil {
			h.logger.Error("failed to clear user cache", "error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
} 
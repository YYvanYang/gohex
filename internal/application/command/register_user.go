package command

import (
	"context"

	"github.com/gohex/gohex/internal/domain/aggregate"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/internal/application/port"
	"github.com/gohex/gohex/pkg/errors"
)

type RegisterUserCommand struct {
	Email    string
	Password string
	Name     string
	Bio      string
}

type RegisterUserResult struct {
	ID string
}

type RegisterUserHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	eventBus   port.EventBus
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func NewRegisterUserHandler(
	userRepo port.UserRepository,
	eventStore port.EventStore,
	eventBus port.EventBus,
	uow port.UnitOfWork,
	logger Logger,
	metrics MetricsReporter,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:   userRepo,
		eventStore: eventStore,
		eventBus:   eventBus,
		uow:        uow,
		logger:     logger,
		metrics:    metrics,
	}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (RegisterUserResult, error) {
	span, ctx := tracer.StartSpan(ctx, "RegisterUserHandler.Handle")
	defer span.End()

	timer := h.metrics.StartTimer("register_user_duration")
	defer timer.Stop()

	var result RegisterUserResult

	err := h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 创建值对象
		email, err := vo.NewEmail(cmd.Email)
		if err != nil {
			return err
		}

		password, err := vo.NewPassword(cmd.Password)
		if err != nil {
			return err
		}

		profile, err := vo.NewUserProfile(cmd.Name, cmd.Bio)
		if err != nil {
			return err
		}

		// 2. 检查用户是否已存在
		exists, err := h.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailAlreadyExists
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

		// 6. 发布事件
		for _, event := range user.Events() {
			if err := h.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}

		result.ID = user.ID()
		return nil
	})

	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		h.metrics.IncrementCounter("register_user_failure")
		return result, err
	}

	h.metrics.IncrementCounter("register_user_success")
	return result, nil
} 
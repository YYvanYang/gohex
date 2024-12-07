package command

import (
	"context"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/errors"
)

// ChangeUserStatusCommand 修改用户状态命令
type ChangeUserStatusCommand struct {
	UserID string
	Status string
}

type ChangeUserStatusHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *ChangeUserStatusHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	statusCmd := cmd.(*ChangeUserStatusCommand)

	return nil, h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 获取用户
		user, err := h.userRepo.FindByID(ctx, statusCmd.UserID)
		if err != nil {
			return err
		}

		// 2. 转换状态值对象
		status := vo.UserStatus(statusCmd.Status)
		if !status.IsValid() {
			return errors.NewValidationError("invalid status")
		}

		// 3. 修改状态
		if err := user.ChangeStatus(status); err != nil {
			return err
		}

		// 4. 保存用户
		if err := h.userRepo.Update(ctx, user); err != nil {
			return err
		}

		// 5. 保存事件
		return h.eventStore.SaveEvents(ctx, user.ID(), user.Events(), user.Version())
	})
} 
package command

import (
	"context"
	"github.com/your-org/your-project/internal/domain/vo"
	"github.com/your-org/your-project/pkg/errors"
)

// AssignRoleCommand 分配角色命令
type AssignRoleCommand struct {
	UserID string
	Role   string
}

type AssignRoleHandler struct {
	userRepo   port.UserRepository
	eventStore port.EventStore
	uow        port.UnitOfWork
	logger     Logger
	metrics    MetricsReporter
}

func (h *AssignRoleHandler) Handle(ctx context.Context, cmd interface{}) (interface{}, error) {
	assignCmd := cmd.(*AssignRoleCommand)

	return nil, h.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 获取用户
		user, err := h.userRepo.FindByID(ctx, assignCmd.UserID)
		if err != nil {
			return err
		}

		// 2. 转换角色值对象
		role := vo.UserRole(assignCmd.Role)
		if !role.IsValid() {
			return errors.NewValidationError("invalid role")
		}

		// 3. 分配角色
		if err := user.AssignRole(role); err != nil {
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
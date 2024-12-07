package event

import (
	"time"
	"github.com/your-org/your-project/internal/domain/vo"
)

const (
	UserCreated       = "user.created"
	UserProfileUpdated = "user.profile_updated"
	PasswordChanged   = "user.password_changed"
	RoleAssigned     = "user.role_assigned"
	UserStatusChanged = "user.status_changed"
	UserDeactivated   = "user.deactivated"
	UserLoggedIn      = "user.logged_in"
	UserLocked        = "user.locked"
	UserUnlocked      = "user.unlocked"
	RoleRevoked       = "user.role_revoked"
)

type UserCreatedEvent struct {
	BaseEvent
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserCreatedEvent(userID string, email string, name string) Event {
	return &UserCreatedEvent{
		BaseEvent: NewBaseEvent(userID, UserCreated),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

type UserProfileUpdatedEvent struct {
	BaseEvent
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserProfileUpdatedEvent(userID string, name string, bio string) Event {
	return &UserProfileUpdatedEvent{
		BaseEvent: NewBaseEvent(userID, UserProfileUpdated),
		Name:      name,
		Bio:       bio,
		UpdatedAt: time.Now(),
	}
}

type UserStatusChangedEvent struct {
	BaseEvent
	OldStatus vo.UserStatus `json:"old_status"`
	NewStatus vo.UserStatus `json:"new_status"`
	ChangedAt time.Time     `json:"changed_at"`
}

func NewUserStatusChangedEvent(userID string, oldStatus, newStatus vo.UserStatus) Event {
	return &UserStatusChangedEvent{
		BaseEvent:  NewBaseEvent(userID, UserStatusChanged),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		ChangedAt:  time.Now(),
	}
}

type UserRoleAssignedEvent struct {
	BaseEvent
	Role      vo.UserRole `json:"role"`
	AssignedAt time.Time  `json:"assigned_at"`
}

func NewUserRoleAssignedEvent(userID string, role vo.UserRole) Event {
	return &UserRoleAssignedEvent{
		BaseEvent:  NewBaseEvent(userID, RoleAssigned),
		Role:      role,
		AssignedAt: time.Now(),
	}
}

type UserLoggedInEvent struct {
	BaseEvent
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	LoginAt   time.Time `json:"login_at"`
}

func NewUserLoggedInEvent(userID string, ip string, userAgent string) Event {
	return &UserLoggedInEvent{
		BaseEvent: NewBaseEvent(userID, UserLoggedIn),
		IP:        ip,
		UserAgent: userAgent,
		LoginAt:   time.Now(),
	}
}

type UserLockedEvent struct {
	BaseEvent
	Reason    string    `json:"reason"`
	LockedAt  time.Time `json:"locked_at"`
}

func NewUserLockedEvent(userID string, reason string) Event {
	return &UserLockedEvent{
		BaseEvent: NewBaseEvent(userID, UserLocked),
		Reason:    reason,
		LockedAt:  time.Now(),
	}
}

type UserUnlockedEvent struct {
	BaseEvent
	UnlockedAt time.Time `json:"unlocked_at"`
}

func NewUserUnlockedEvent(userID string) Event {
	return &UserUnlockedEvent{
		BaseEvent:  NewBaseEvent(userID, UserUnlocked),
		UnlockedAt: time.Now(),
	}
}

type UserRoleRevokedEvent struct {
	BaseEvent
	Role      vo.UserRole `json:"role"`
	RevokedAt time.Time   `json:"revoked_at"`
}

func NewUserRoleRevokedEvent(userID string, role vo.UserRole) Event {
	return &UserRoleRevokedEvent{
		BaseEvent: NewBaseEvent(userID, RoleRevoked),
		Role:      role,
		RevokedAt: time.Now(),
	}
}

// 其他事件类型的实现... 
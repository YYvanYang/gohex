package aggregate

import (
	"github.com/google/uuid"
	"github.com/gohex/gohex/internal/domain/event"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/errors"
	"time"
)

type User struct {
	*BaseAggregate
	email     vo.Email
	password  vo.Password
	profile   vo.UserProfile
	status    vo.UserStatus
	roles     []vo.UserRole
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email vo.Email, password vo.Password, profile vo.UserProfile) (*User, error) {
	user := &User{
		BaseAggregate: NewBaseAggregate(uuid.New().String()),
		email:        email,
		password:     password,
		profile:      profile,
		status:       vo.StatusActive,
		roles:        []vo.UserRole{vo.RoleUser},
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}

	user.AddEvent(event.NewUserCreatedEvent(
		user.ID(),
		email.String(),
		profile.Name(),
	))

	return user, nil
}

// Getters
func (u *User) Email() vo.Email { return u.email }
func (u *User) Password() vo.Password { return u.password }
func (u *User) Profile() vo.UserProfile { return u.profile }
func (u *User) Status() vo.UserStatus { return u.status }
func (u *User) Roles() []vo.UserRole { return u.roles }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Business Methods
func (u *User) UpdateProfile(profile vo.UserProfile) error {
	if profile.IsEmpty() {
		return errors.ErrInvalidProfile
	}

	u.profile = profile
	u.updatedAt = time.Now()

	u.AddEvent(event.NewUserProfileUpdatedEvent(
		u.ID(),
		profile.Name(),
		profile.Bio(),
	))

	return nil
}

func (u *User) ChangePassword(current, new vo.Password) error {
	if err := u.password.Compare(current.Hash()); err != nil {
		return errors.ErrInvalidPassword
	}

	u.password = new
	u.updatedAt = time.Now()

	u.AddEvent(event.NewPasswordChangedEvent(u.ID()))
	return nil
}

func (u *User) ResetPassword(new vo.Password) error {
	u.password = new
	u.updatedAt = time.Now()

	u.AddEvent(event.NewPasswordResetEvent(u.ID()))
	return nil
}

func (u *User) ChangeStatus(status vo.UserStatus) error {
	if !status.IsValid() {
		return errors.ErrInvalidStatus
	}

	if u.status == status {
		return nil
	}

	oldStatus := u.status
	u.status = status
	u.updatedAt = time.Now()

	u.AddEvent(event.NewUserStatusChangedEvent(
		u.ID(),
		oldStatus,
		status,
	))

	return nil
}

func (u *User) AssignRole(role vo.UserRole) error {
	if !role.IsValid() {
		return errors.ErrInvalidRole
	}

	// 检查角色是否已存在
	for _, r := range u.roles {
		if r == role {
			return errors.ErrRoleAlreadyAssigned
		}
	}

	u.roles = append(u.roles, role)
	u.updatedAt = time.Now()

	u.AddEvent(event.NewUserRoleAssignedEvent(u.ID(), role))
	return nil
}

func (u *User) RevokeRole(role vo.UserRole) error {
	if !role.IsValid() {
		return errors.ErrInvalidRole
	}

	// 不能移除最后一个角色
	if len(u.roles) == 1 {
		return errors.ErrCannotRevokeLastRole
	}

	// 移除角色
	var newRoles []vo.UserRole
	found := false
	for _, r := range u.roles {
		if r != role {
			newRoles = append(newRoles, r)
		} else {
			found = true
		}
	}

	if !found {
		return errors.ErrRoleNotFound
	}

	u.roles = newRoles
	u.updatedAt = time.Now()

	u.AddEvent(event.NewUserRoleRevokedEvent(u.ID(), role))
	return nil
}

func (u *User) ValidatePassword(plaintext string) error {
	return u.password.Compare(plaintext)
}

func (u *User) Lock() error {
	if err := u.ChangeStatus(vo.StatusSuspended); err != nil {
		return err
	}
	u.AddEvent(event.NewUserLockedEvent(u.ID()))
	return nil
}

func (u *User) Unlock() error {
	return u.ChangeStatus(vo.StatusActive)
}

func (u *User) IsActive() bool {
	return u.status.IsActive()
}

func (u *User) HasRole(role vo.UserRole) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.roles {
		if role.HasPermission(permission) {
			return true
		}
	}
	return false
}

func (u *User) RoleStrings() []string {
	roles := make([]string, len(u.roles))
	for i, role := range u.roles {
		roles[i] = role.String()
	}
	return roles
}

func (u *User) RecordLogin(ip string, userAgent string) {
	u.AddEvent(event.NewUserLoggedInEvent(
		u.ID(),
		ip,
		userAgent,
	))
}
 
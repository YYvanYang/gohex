package vo

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusDeleted   UserStatus = "deleted"
)

var validStatuses = map[UserStatus]bool{
	StatusActive:    true,
	StatusInactive:  true,
	StatusSuspended: true,
	StatusDeleted:   true,
}

func (s UserStatus) IsValid() bool {
	return validStatuses[s]
}

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsActive() bool {
	return s == StatusActive
}

func (s UserStatus) CanBeActivated() bool {
	return s == StatusInactive || s == StatusSuspended
}

func (s UserStatus) CanBeDeactivated() bool {
	return s == StatusActive
} 
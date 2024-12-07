package vo

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleMod   UserRole = "moderator"
)

var validRoles = map[UserRole]bool{
	RoleUser:  true,
	RoleAdmin: true,
	RoleMod:   true,
}

func (r UserRole) IsValid() bool {
	return validRoles[r]
}

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) HasPermission(permission string) bool {
	switch r {
	case RoleAdmin:
		return true
	case RoleMod:
		return isModeratorPermission(permission)
	case RoleUser:
		return isUserPermission(permission)
	default:
		return false
	}
}

func isModeratorPermission(permission string) bool {
	moderatorPermissions := map[string]bool{
		"users.view":   true,
		"users.update": true,
		"content.moderate": true,
	}
	return moderatorPermissions[permission]
}

func isUserPermission(permission string) bool {
	userPermissions := map[string]bool{
		"profile.view":   true,
		"profile.update": true,
	}
	return userPermissions[permission]
} 
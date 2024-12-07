package dto

import "time"

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Status    string    `json:"status"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListDTO 用户列表数据传输对象
type UserListDTO struct {
	Total int64      `json:"total"`
	Items []UserDTO  `json:"items"`
}

// CreateUserDTO 创建用户请求
type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"max=500"`
}

// UpdateUserDTO 更新用户请求
type UpdateUserDTO struct {
	Name   string `json:"name" validate:"required"`
	Bio    string `json:"bio" validate:"max=500"`
	Avatar string `json:"avatar" validate:"omitempty,url"`
}

// ChangePasswordDTO 修改密码请求
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UserProfileDTO 用户配置数据传输对象
type UserProfileDTO struct {
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"max=500"`
	Avatar   string `json:"avatar,omitempty" validate:"omitempty,url"`
	Location string `json:"location,omitempty"`
	Website  string `json:"website,omitempty" validate:"omitempty,url"`
}
 
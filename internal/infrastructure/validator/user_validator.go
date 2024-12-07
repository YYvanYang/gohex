package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/your-org/your-project/internal/domain/vo"
)

type UserValidator struct {
	validate *validator.Validate
}

func NewUserValidator() *UserValidator {
	v := validator.New()

	// 注册自定义验证
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("user_status", validateUserStatus)
	v.RegisterValidation("user_role", validateUserRole)

	return &UserValidator{validate: v}
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := vo.NewEmail(email)
	return err == nil
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	// 其他密码复杂度验证...
	return true
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	// 其他用户名规则验证...
	return true
}

func validateUserStatus(fl validator.FieldLevel) bool {
	status := vo.UserStatus(fl.Field().String())
	return status.IsValid()
}

func validateUserRole(fl validator.FieldLevel) bool {
	role := vo.UserRole(fl.Field().String())
	return role.IsValid()
} 
package dto

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func init() {
    // 注册自定义验证规则
    validate.RegisterValidation("password", validatePassword)
    validate.RegisterValidation("username", validateUsername)
}

func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    // 密码规则：至少8位，包含大小写字母和数字
    // TODO: 实现密码复杂度验证
    return len(password) >= 8
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // 用户名规则：3-20位字母数字下划线
    // TODO: 实现用户名格式验证
    return len(username) >= 3 && len(username) <= 20
}

func Validate(i interface{}) error {
    return validate.Struct(i)
} 
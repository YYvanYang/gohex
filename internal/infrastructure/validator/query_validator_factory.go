package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/your-org/your-project/internal/application/port/input/query"
)

type ValidatorFactory interface {
	CreateQueryValidator() query.Validator
	CreateCommandValidator() command.Validator
}

type validatorFactory struct {
	validate *validator.Validate
	logger   Logger
}

func NewValidatorFactory(logger Logger) ValidatorFactory {
	v := validator.New()
	
	// 注册自定义验证规则
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	
	return &validatorFactory{
		validate: v,
		logger:   logger,
	}
}

func (f *validatorFactory) CreateQueryValidator() query.Validator {
	return NewQueryValidator(f.validate, f.logger)
}

func (f *validatorFactory) CreateCommandValidator() command.Validator {
	return NewCommandValidator(f.validate, f.logger)
} 
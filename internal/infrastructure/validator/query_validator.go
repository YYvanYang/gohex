package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/your-org/your-project/internal/application/port/input/query"
)

type queryValidator struct {
	validate *validator.Validate
}

func NewQueryValidator() query.Validator {
	return &queryValidator{
		validate: validator.New(),
	}
}

func (v *queryValidator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
} 
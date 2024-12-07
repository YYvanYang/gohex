package errors

var (
	ErrInvalidEmail = &AppError{
		Code:    ErrCodeValidation,
		Message: "invalid email format",
	}
	
	ErrInvalidPassword = &AppError{
		Code:    ErrCodeValidation,
		Message: "invalid password format",
	}
) 
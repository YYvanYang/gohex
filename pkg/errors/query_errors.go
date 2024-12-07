package errors

var (
    ErrUserNotFound = &AppError{
        Code:    ErrCodeNotFound,
        Message: "user not found",
    }
) 
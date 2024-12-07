package query

import "github.com/your-org/your-project/pkg/errors"

type QueryError struct {
    Code    string
    Message string
    Cause   error
}

func (e *QueryError) Error() string {
    if e.Cause != nil {
        return e.Message + ": " + e.Cause.Error()
    }
    return e.Message
}

func NewQueryError(code string, message string, cause error) error {
    return &QueryError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
}

var (
    ErrQueryValidation = func(cause error) error {
        return NewQueryError("VALIDATION_ERROR", "Query validation failed", cause)
    }
    
    ErrQueryExecution = func(cause error) error {
        return NewQueryError("EXECUTION_ERROR", "Query execution failed", cause)
    }
    
    ErrQueryTimeout = func(cause error) error {
        return NewQueryError("TIMEOUT_ERROR", "Query timed out", cause)
    }
) 
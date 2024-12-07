package query

import "fmt"

var (
    ErrQueryNotFound    = NewQueryError("QUERY_NOT_FOUND", "Query handler not found")
    ErrQueryValidation  = NewQueryError("QUERY_VALIDATION", "Query validation failed")
    ErrQueryExecution   = NewQueryError("QUERY_EXECUTION", "Query execution failed")
    ErrQueryTimeout     = NewQueryError("QUERY_TIMEOUT", "Query execution timed out")
    ErrQueryCacheMiss   = NewQueryError("CACHE_MISS", "Cache miss")
    ErrQueryPermission  = NewQueryError("PERMISSION_DENIED", "Permission denied")
)

type QueryError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}

func (e *QueryError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewQueryError(code string, message string) *QueryError {
    return &QueryError{
        Code:    code,
        Message: message,
        Details: make(map[string]interface{}),
    }
}

func (e *QueryError) WithCause(cause error) *QueryError {
    e.Cause = cause
    return e
}

func (e *QueryError) WithDetail(key string, value interface{}) *QueryError {
    e.Details[key] = value
    return e
} 
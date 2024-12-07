package errors

// 系统错误码
const (
    ErrCodeSystem           = "GOHEX_SYSTEM_ERROR"
    ErrCodeDatabase        = "GOHEX_DATABASE_ERROR"
    ErrCodeCache           = "GOHEX_CACHE_ERROR"
    ErrCodeNetwork         = "GOHEX_NETWORK_ERROR"
    ErrCodeConfiguration   = "GOHEX_CONFIG_ERROR"
    ErrCodeSerialization   = "GOHEX_SERIALIZATION_ERROR"
)

// 业务错误码
const (
    ErrCodeBusiness        = "GOHEX_BUSINESS_ERROR"
    ErrCodeInvalidState    = "GOHEX_INVALID_STATE"
    ErrCodeRateLimit       = "GOHEX_RATE_LIMIT"
    ErrCodeResourceLocked  = "GOHEX_RESOURCE_LOCKED"
    ErrCodeDependencyFailed = "GOHEX_DEPENDENCY_FAILED"
)

// 安全错误码
const (
    ErrCodeSecurity        = "SECURITY_ERROR"
    ErrCodeAccessDenied    = "ACCESS_DENIED"
    ErrCodeInvalidSession  = "INVALID_SESSION"
    ErrCodePasswordExpired = "PASSWORD_EXPIRED"
    ErrCodeAccountLocked   = "ACCOUNT_LOCKED"
)

// 验证错误码
const (
    ErrCodeValidationFailed = "VALIDATION_FAILED"
    ErrCodeInvalidFormat    = "INVALID_FORMAT"
    ErrCodeMissingField     = "MISSING_FIELD"
    ErrCodeInvalidValue     = "INVALID_VALUE"
    ErrCodeDuplicateValue   = "DUPLICATE_VALUE"
    ErrCodeValidation   ErrorCode = "GOHEX_VALIDATION_ERROR"
    ErrCodeNotFound     ErrorCode = "GOHEX_NOT_FOUND"
    ErrCodeConflict     ErrorCode = "GOHEX_CONFLICT"
    ErrCodeUnauthorized ErrorCode = "GOHEX_UNAUTHORIZED"
    ErrCodeInternal     ErrorCode = "GOHEX_INTERNAL_ERROR"
) 
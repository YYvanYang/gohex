package output

import (
    "context"
    "time"
)

type AuthService interface {
    // GenerateToken 生成令牌
    GenerateToken(user *aggregate.User) (string, time.Time, error)
    // ValidateToken 验证令牌
    ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
    // RevokeToken 吊销令牌
    RevokeToken(ctx context.Context, token string) error
    // RefreshToken 刷新令牌
    RefreshToken(ctx context.Context, refreshToken string) (string, time.Time, error)
    // GetTokenInfo 获取令牌信息
    GetTokenInfo(ctx context.Context, token string) (*TokenInfo, error)
    // IsTokenRevoked 检查令牌是否已吊销
    IsTokenRevoked(ctx context.Context, token string) bool
}

type TokenClaims struct {
    UserID    string
    Email     string
    Roles     []string
    ExpiresAt time.Time
}

type TokenInfo struct {
    TokenClaims
    IssuedAt  time.Time
    NotBefore time.Time
    Issuer    string
    Subject   string
    Audience  []string
} 
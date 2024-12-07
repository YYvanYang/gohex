package output

import (
	"context"
	"time"
	"github.com/your-org/your-project/internal/domain/aggregate"
)

// TokenClaims 令牌声明
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	ExpiresAt time.Time `json:"exp"`
}

// TokenService 令牌服务接口
type TokenService interface {
	GenerateToken(user *aggregate.User) (string, time.Time, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	RevokeToken(ctx context.Context, token string) error
} 
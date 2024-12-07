package jwt

import (
	"context"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gohex/gohex/internal/domain/aggregate"
	"github.com/gohex/gohex/internal/application/port"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/pkg/tracer"
)

type Config struct {
	SecretKey     string
	TokenDuration time.Duration
}

type jwtTokenService struct {
	config  Config
	cache   port.Cache // 用于存储已吊销的令牌
	logger  Logger
	metrics MetricsReporter
}

func NewJWTTokenService(
	config Config,
	cache port.Cache,
	logger Logger,
	metrics MetricsReporter,
) port.TokenService {
	return &jwtTokenService{
		config:  config,
		cache:   cache,
		logger:  logger,
		metrics: metrics,
	}
}

func (s *jwtTokenService) GenerateToken(user *aggregate.User) (string, time.Time, error) {
	timer := s.metrics.StartTimer("token_generation_duration")
	defer timer.Stop()

	expiresAt := time.Now().Add(s.config.TokenDuration)

	claims := jwt.MapClaims{
		"user_id": user.ID(),
		"email":   user.Email().String(),
		"roles":   user.RoleStrings(),
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		s.logger.Error("failed to sign token", "error", err)
		s.metrics.IncrementCounter("token_generation_failure")
		return "", time.Time{}, err
	}

	s.metrics.IncrementCounter("token_generation_success")
	return signedToken, expiresAt, nil
}

func (s *jwtTokenService) ValidateToken(ctx context.Context, tokenString string) (*port.TokenClaims, error) {
	timer := s.metrics.StartTimer("token_validation_duration")
	defer timer.Stop()

	// 1. 检查令牌是否被吊销
	if s.isTokenRevoked(ctx, tokenString) {
		return nil, errors.ErrTokenRevoked
	}

	// 2. 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		s.metrics.IncrementCounter("token_validation_failure")
		return nil, errors.ErrInvalidToken
	}

	// 3. 验证声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.metrics.IncrementCounter("token_validation_failure")
		return nil, errors.ErrInvalidToken
	}

	// 4. 转换声明
	expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)
	roles := make([]string, len(claims["roles"].([]interface{})))
	for i, role := range claims["roles"].([]interface{}) {
		roles[i] = role.(string)
	}

	s.metrics.IncrementCounter("token_validation_success")
	return &port.TokenClaims{
		UserID:    claims["user_id"].(string),
		Email:     claims["email"].(string),
		Roles:     roles,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *jwtTokenService) RevokeToken(ctx context.Context, token string) error {
	timer := s.metrics.StartTimer("token_revocation_duration")
	defer timer.Stop()

	// 将令牌加入黑名单
	key := "revoked_token:" + token
	if err := s.cache.Set(ctx, key, true, s.config.TokenDuration); err != nil {
		s.logger.Error("failed to revoke token", "error", err)
		s.metrics.IncrementCounter("token_revocation_failure")
		return err
	}

	s.metrics.IncrementCounter("token_revocation_success")
	return nil
}

func (s *jwtTokenService) isTokenRevoked(ctx context.Context, token string) bool {
	key := "revoked_token:" + token
	revoked, _ := s.cache.Get(ctx, key)
	return revoked != nil
} 
package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:uuid;not null;index"                   json:"user_id"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex;not null"                    json:"-"`
	UserAgent string     `gorm:"column:user_agent"                                         json:"user_agent"`
	IPAddress string     `gorm:"column:ip_address"                                         json:"ip_address"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"                                json:"expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"                                         json:"revoked_at,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"                          json:"created_at"`
}

func (RefreshSession) TableName() string { return "refresh_sessions" }

func (s *RefreshSession) IsValid() bool {
	return s.RevokedAt == nil && time.Now().Before(s.ExpiresAt)
}

type RefreshSessionRepository interface {
	Create(ctx context.Context, s *RefreshSession) error
	FindByTokenHash(ctx context.Context, hash string) (*RefreshSession, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

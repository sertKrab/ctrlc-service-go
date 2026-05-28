package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	domainauth "git.trovefin.com/poc/ctrlc-service-go/internal/domain/auth"
)

type refreshSessionRepo struct{ db *gorm.DB }

func NewRefreshSessionRepo(db *gorm.DB) domainauth.RefreshSessionRepository {
	return &refreshSessionRepo{db: db}
}

func (r *refreshSessionRepo) Create(ctx context.Context, s *domainauth.RefreshSession) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *refreshSessionRepo) FindByTokenHash(ctx context.Context, hash string) (*domainauth.RefreshSession, error) {
	var s domainauth.RefreshSession
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *refreshSessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domainauth.RefreshSession{}).
		Where("id = ?", id).Update("revoked_at", now).Error
}

func (r *refreshSessionRepo) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&domainauth.RefreshSession{}).Error
}

package repository

import (
	"context"

	"gorm.io/gorm"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/audit"
)

type auditRepo struct{ db *gorm.DB }

func NewAuditRepo(db *gorm.DB) audit.Repository { return &auditRepo{db: db} }

func (r *auditRepo) Log(ctx context.Context, entry audit.Entry) error {
	log := audit.AuditLog{
		UserID:     entry.UserID,
		Action:     entry.Action,
		EntityType: entry.EntityType,
		EntityID:   entry.EntityID,
		Changes:    entry.Changes,
		IPAddress:  entry.IPAddress,
		UserAgent:  entry.UserAgent,
	}
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *auditRepo) List(ctx context.Context, entityType, entityID string, page, pageSize int) ([]audit.AuditLog, int64, error) {
	var logs []audit.AuditLog
	var total int64
	q := r.db.WithContext(ctx).Model(&audit.AuditLog{})
	if entityType != "" {
		q = q.Where("entity_type = ?", entityType)
	}
	if entityID != "" {
		q = q.Where("entity_id = ?", entityID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

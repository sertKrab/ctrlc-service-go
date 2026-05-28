package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     *uuid.UUID     `gorm:"column:user_id;type:uuid;index"                           json:"user_id,omitempty"`
	Action     string         `gorm:"column:action;not null"                                   json:"action"`
	EntityType string         `gorm:"column:entity_type;not null;index"                        json:"entity_type"`
	EntityID   string         `gorm:"column:entity_id;index"                                   json:"entity_id"`
	Changes    datatypes.JSON `gorm:"column:changes;type:jsonb"                                json:"changes,omitempty"`
	IPAddress  string         `gorm:"column:ip_address"                                        json:"ip_address"`
	UserAgent  string         `gorm:"column:user_agent"                                        json:"user_agent"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"                         json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }

type Entry struct {
	UserID     *uuid.UUID
	Action     string
	EntityType string
	EntityID   string
	Changes    datatypes.JSON
	IPAddress  string
	UserAgent  string
}

type Repository interface {
	Log(ctx context.Context, entry Entry) error
	List(ctx context.Context, entityType, entityID string, page, pageSize int) ([]AuditLog, int64, error)
}

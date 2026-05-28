package permission

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Resource    string    `gorm:"column:resource;not null"                                  json:"resource"`
	Action      string    `gorm:"column:action;not null"                                    json:"action"`
	Description string    `gorm:"column:description"                                        json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"                          json:"created_at"`
}

func (Permission) TableName() string { return "permissions" }

type Repository interface {
	FindAll(ctx context.Context) ([]Permission, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	FindByResourceAction(ctx context.Context, resource, action string) (*Permission, error)
}

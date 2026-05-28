package role

import (
	"context"
	"time"

	"github.com/google/uuid"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/permission"
)

type Role struct {
	ID          uuid.UUID              `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string                 `gorm:"column:name;uniqueIndex;not null"                          json:"name"`
	Description string                 `gorm:"column:description"                                        json:"description"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions"                               json:"permissions,omitempty"`
	CreatedAt   time.Time              `gorm:"column:created_at;autoCreateTime"                          json:"created_at"`
	UpdatedAt   time.Time              `gorm:"column:updated_at;autoUpdateTime"                          json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

type Repository interface {
	FindAll(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
}

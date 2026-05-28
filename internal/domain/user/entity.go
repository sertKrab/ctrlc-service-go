package user

import (
	"context"

	"github.com/google/uuid"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/role"
)

type User struct {
	domain.Base
	Email        string    `gorm:"column:email;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null"     json:"-"`
	FirstName    string    `gorm:"column:first_name;not null"        json:"first_name"`
	LastName     string    `gorm:"column:last_name;not null"         json:"last_name"`
	IsActive     bool      `gorm:"column:is_active;default:true"     json:"is_active"`
	RoleID       uuid.UUID `gorm:"column:role_id;type:uuid;not null" json:"role_id"`
	Role         role.Role `gorm:"foreignKey:RoleID"                 json:"role,omitempty"`
}

func (User) TableName() string { return "users" }

type ListRequest struct {
	Page     int
	PageSize int
	Search   string
	RoleID   *uuid.UUID
	IsActive *bool
}

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req ListRequest) ([]User, int64, error)
}

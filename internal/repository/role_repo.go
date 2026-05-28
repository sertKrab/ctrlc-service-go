package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/role"
)

type roleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) role.Repository { return &roleRepo{db: db} }

func (r *roleRepo) FindAll(ctx context.Context) ([]role.Role, error) {
	var roles []role.Role
	return roles, r.db.WithContext(ctx).Preload("Permissions").Find(&roles).Error
}

func (r *roleRepo) FindByID(ctx context.Context, id uuid.UUID) (*role.Role, error) {
	var ro role.Role
	return &ro, r.db.WithContext(ctx).Preload("Permissions").First(&ro, "id = ?", id).Error
}

func (r *roleRepo) FindByName(ctx context.Context, name string) (*role.Role, error) {
	var ro role.Role
	return &ro, r.db.WithContext(ctx).Where("name = ?", name).First(&ro).Error
}

func (r *roleRepo) Create(ctx context.Context, ro *role.Role) error {
	return r.db.WithContext(ctx).Create(ro).Error
}

func (r *roleRepo) Update(ctx context.Context, ro *role.Role) error {
	return r.db.WithContext(ctx).Save(ro).Error
}

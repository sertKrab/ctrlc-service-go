package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/permission"
)

type permissionRepo struct{ db *gorm.DB }

func NewPermissionRepo(db *gorm.DB) permission.Repository { return &permissionRepo{db: db} }

func (r *permissionRepo) FindAll(ctx context.Context) ([]permission.Permission, error) {
	var perms []permission.Permission
	return perms, r.db.WithContext(ctx).Find(&perms).Error
}

func (r *permissionRepo) FindByID(ctx context.Context, id uuid.UUID) (*permission.Permission, error) {
	var p permission.Permission
	return &p, r.db.WithContext(ctx).First(&p, "id = ?", id).Error
}

func (r *permissionRepo) FindByResourceAction(ctx context.Context, resource, action string) (*permission.Permission, error) {
	var p permission.Permission
	return &p, r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&p).Error
}

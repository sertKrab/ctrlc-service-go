package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type userRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) user.Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Preload("Role.Permissions").
		Where("id = ? AND deleted_at IS NULL", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Preload("Role.Permissions").
		Where("email = ? AND deleted_at IS NULL", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error
}

func (r *userRepo) List(ctx context.Context, req user.ListRequest) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	q := r.db.WithContext(ctx).Model(&user.User{}).Preload("Role")
	if req.Search != "" {
		q = q.Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.RoleID != nil {
		q = q.Where("role_id = ?", *req.RoleID)
	}
	if req.IsActive != nil {
		q = q.Where("is_active = ?", *req.IsActive)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

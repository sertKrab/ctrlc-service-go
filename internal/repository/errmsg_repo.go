package repository

import (
	"context"

	"gorm.io/gorm"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
)

type errMsgRepo struct{ db *gorm.DB }

func NewErrMsgRepo(db *gorm.DB) errmsg.Repository { return &errMsgRepo{db: db} }

func (r *errMsgRepo) FindAll(ctx context.Context) ([]errmsg.ErrorMessage, error) {
	var msgs []errmsg.ErrorMessage
	return msgs, r.db.WithContext(ctx).Find(&msgs).Error
}

func (r *errMsgRepo) FindByCode(ctx context.Context, code string) (*errmsg.ErrorMessage, error) {
	var msg errmsg.ErrorMessage
	return &msg, r.db.WithContext(ctx).Where("code = ?", code).First(&msg).Error
}

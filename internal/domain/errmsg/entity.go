package errmsg

import (
	"context"
	"time"
)

type ErrorMessage struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"column:code;uniqueIndex;not null"   json:"code"`
	LocaleTH  string    `gorm:"column:locale_th;not null"          json:"locale_th"`
	LocaleEN  string    `gorm:"column:locale_en;not null"          json:"locale_en"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"   json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"   json:"updated_at"`
}

func (ErrorMessage) TableName() string { return "error_messages" }

type Repository interface {
	FindAll(ctx context.Context) ([]ErrorMessage, error)
	FindByCode(ctx context.Context, code string) (*ErrorMessage, error)
}

// Resolver is implemented by the in-memory cache (Phase 6)
type Resolver interface {
	Error(code string) error
	Message(code string) string
}

// AppError is the localized error type returned by all usecases
type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string { return e.Message }

package userdto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email     string    `json:"email"      binding:"required,email"`
	Password  string    `json:"password"   binding:"required,min=8"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name"  binding:"required"`
	RoleID    uuid.UUID `json:"role_id"    binding:"required"`
}

type UpdateUserRequest struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	RoleID    *uuid.UUID `json:"role_id"`
	IsActive  *bool      `json:"is_active"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}

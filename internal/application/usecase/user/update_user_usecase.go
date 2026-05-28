package userusecase

import (
	"context"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type UpdateUserUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewUpdateUserUseCase(userRepo user.Repository, resolver errmsg.Resolver) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo, resolver}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id interface{}, req interface{}) (interface{}, error) {
	// TODO: implement
	return nil, nil
}

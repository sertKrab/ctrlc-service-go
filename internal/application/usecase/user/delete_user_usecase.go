package userusecase

import (
	"context"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type DeleteUserUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewDeleteUserUseCase(userRepo user.Repository, resolver errmsg.Resolver) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepo, resolver}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id interface{}) error {
	// TODO: implement
	return nil
}

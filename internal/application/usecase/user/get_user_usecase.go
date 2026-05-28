package userusecase

import (
	"context"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type GetUserUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewGetUserUseCase(userRepo user.Repository, resolver errmsg.Resolver) *GetUserUseCase {
	return &GetUserUseCase{userRepo, resolver}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id interface{}) (interface{}, error) {
	// TODO: implement
	return nil, nil
}

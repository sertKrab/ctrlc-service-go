package userusecase

import (
	"context"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type CreateUserUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewCreateUserUseCase(userRepo user.Repository, resolver errmsg.Resolver) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo, resolver}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: implement
	return nil, nil
}

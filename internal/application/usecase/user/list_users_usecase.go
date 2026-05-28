package userusecase

import (
	"context"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type ListUsersUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewListUsersUseCase(userRepo user.Repository, resolver errmsg.Resolver) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo, resolver}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: implement
	return nil, nil
}

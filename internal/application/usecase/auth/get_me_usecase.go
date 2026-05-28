package authusecase

import (
	"context"

	"github.com/google/uuid"
	authdto "git.trovefin.com/poc/ctrlc-service-go/internal/application/dto/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
)

type GetMeUseCase struct {
	userRepo user.Repository
	resolver errmsg.Resolver
}

func NewGetMeUseCase(userRepo user.Repository, resolver errmsg.Resolver) *GetMeUseCase {
	return &GetMeUseCase{userRepo, resolver}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, userID uuid.UUID) (*authdto.UserInfo, error) {
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, uc.resolver.Error("USER_NOT_FOUND")
	}
	return &authdto.UserInfo{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role.Name,
	}, nil
}

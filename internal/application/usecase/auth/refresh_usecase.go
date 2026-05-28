package authusecase

import (
	"context"

	authdto "git.trovefin.com/poc/ctrlc-service-go/internal/application/dto/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	domainauth "git.trovefin.com/poc/ctrlc-service-go/internal/domain/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
	"git.trovefin.com/poc/ctrlc-service-go/pkg/jwtutil"
)

type RefreshUseCase struct {
	sessionRepo domainauth.RefreshSessionRepository
	userRepo    user.Repository
	cfg         *config.Config
	resolver    errmsg.Resolver
}

func NewRefreshUseCase(
	sessionRepo domainauth.RefreshSessionRepository,
	userRepo user.Repository,
	cfg *config.Config,
	resolver errmsg.Resolver,
) *RefreshUseCase {
	return &RefreshUseCase{sessionRepo, userRepo, cfg, resolver}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, rawToken string, meta usecase.RequestMeta) (*authdto.RefreshResponse, error) {
	tokenHash := hashToken(rawToken)

	session, err := uc.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil || session == nil || !session.IsValid() {
		return nil, uc.resolver.Error("AUTH_SESSION_REVOKED")
	}

	u, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil || !u.IsActive {
		return nil, uc.resolver.Error("AUTH_UNAUTHORIZED")
	}

	if err := uc.sessionRepo.Revoke(ctx, session.ID); err != nil {
		return nil, uc.resolver.Error("INTERNAL_ERROR")
	}

	loginUC := &LoginUseCase{
		userRepo:    uc.userRepo,
		sessionRepo: uc.sessionRepo,
		cfg:         uc.cfg,
		resolver:    uc.resolver,
	}
	newRefreshToken, err := loginUC.createSession(ctx, u.ID, meta)
	if err != nil {
		return nil, uc.resolver.Error("INTERNAL_ERROR")
	}

	accessToken, expiresIn, err := jwtutil.GenerateAccessToken(u.ID, u.RoleID, uc.cfg.JWTSecret, uc.cfg.AccessTokenExpiry)
	if err != nil {
		return nil, uc.resolver.Error("INTERNAL_ERROR")
	}

	return &authdto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

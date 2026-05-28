package authusecase

import (
	"context"

	"github.com/google/uuid"
	"git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/audit"
	domainauth "git.trovefin.com/poc/ctrlc-service-go/internal/domain/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
)

type LogoutUseCase struct {
	sessionRepo domainauth.RefreshSessionRepository
	auditRepo   audit.Repository
	resolver    errmsg.Resolver
}

func NewLogoutUseCase(
	sessionRepo domainauth.RefreshSessionRepository,
	auditRepo audit.Repository,
	resolver errmsg.Resolver,
) *LogoutUseCase {
	return &LogoutUseCase{sessionRepo, auditRepo, resolver}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, rawToken string, userID uuid.UUID, meta usecase.RequestMeta) error {
	if rawToken == "" {
		return nil
	}
	tokenHash := hashToken(rawToken)
	session, err := uc.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil || session == nil {
		return nil // already gone — success
	}
	if err := uc.sessionRepo.Revoke(ctx, session.ID); err != nil {
		return uc.resolver.Error("INTERNAL_ERROR")
	}
	go func() {
		_ = uc.auditRepo.Log(context.Background(), audit.Entry{
			UserID: &userID, Action: "LOGOUT",
			EntityType: "user", EntityID: userID.String(),
			IPAddress: meta.IPAddress, UserAgent: meta.UserAgent,
		})
	}()
	return nil
}

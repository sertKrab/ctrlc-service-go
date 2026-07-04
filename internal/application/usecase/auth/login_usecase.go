package authusecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	authdto "git.trovefin.com/poc/ctrlc-service-go/internal/application/dto/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/audit"
	domainauth "git.trovefin.com/poc/ctrlc-service-go/internal/domain/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/logging"
	"git.trovefin.com/poc/ctrlc-service-go/pkg/jwtutil"
)

type LoginUseCase struct {
	userRepo    user.Repository
	sessionRepo domainauth.RefreshSessionRepository
	auditRepo   audit.Repository
	cfg         *config.Config
	resolver    errmsg.Resolver
}

func NewLoginUseCase(
	userRepo user.Repository,
	sessionRepo domainauth.RefreshSessionRepository,
	auditRepo audit.Repository,
	cfg *config.Config,
	resolver errmsg.Resolver,
) *LoginUseCase {
	return &LoginUseCase{userRepo, sessionRepo, auditRepo, cfg, resolver}
}

// Execute demonstrates the four standard log shapes this template's logging
// package supports — use this as the reference pattern for new modules:
//  1. general log:          log.Debug("message", "key", value)
//  2. function-tagged log:  `log` below is tagged with KeyFunc once, reused for every call in this function
//  3. user-scoped log:      log.Warn("message", logging.KeyUserID, u.ID.String())
//  4. masked sensitive log: logging.Mask(req.Email) — email isn't resolved to a user yet, so log it masked
func (uc *LoginUseCase) Execute(ctx context.Context, req authdto.LoginRequest, meta usecase.RequestMeta) (*authdto.LoginResponse, error) {
	log := logging.FromContext(ctx).With(logging.KeyFunc, "LoginUseCase.Execute")
	log.Debug("login attempt", "email", logging.Mask(req.Email))

	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Warn("user not found", "email", logging.Mask(req.Email))
		return nil, uc.resolver.Error("AUTH_INVALID_CREDENTIALS")
	}

	if !u.IsActive {
		log.Warn("user inactive", logging.KeyUserID, u.ID.String())
		return nil, uc.resolver.Error("USER_INACTIVE")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn("invalid password", logging.KeyUserID, u.ID.String())
		return nil, uc.resolver.Error("AUTH_INVALID_CREDENTIALS")
	}

	accessToken, expiresIn, err := jwtutil.GenerateAccessToken(u.ID, u.RoleID, uc.cfg.JWTSecret, uc.cfg.AccessTokenExpiry)
	if err != nil {
		log.Error("access token generation failed", logging.KeyUserID, u.ID.String(), "error", err)
		return nil, uc.resolver.Error("INTERNAL_ERROR")
	}

	refreshToken, err := uc.createSession(ctx, u.ID, meta)
	if err != nil {
		log.Error("session creation failed", logging.KeyUserID, u.ID.String(), "error", err)
		return nil, uc.resolver.Error("INTERNAL_ERROR")
	}

	log.Info("login success", logging.KeyUserID, u.ID.String())

	go func() {
		_ = uc.auditRepo.Log(context.Background(), audit.Entry{
			UserID: &u.ID, Action: "LOGIN",
			EntityType: "user", EntityID: u.ID.String(),
			IPAddress: meta.IPAddress, UserAgent: meta.UserAgent,
		})
	}()

	return &authdto.LoginResponse{
		User: authdto.UserInfo{
			ID: u.ID, Email: u.Email,
			FirstName: u.FirstName, LastName: u.LastName,
			Role: u.Role.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (uc *LoginUseCase) createSession(ctx context.Context, userID uuid.UUID, meta usecase.RequestMeta) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	rawToken := hex.EncodeToString(raw)
	tokenHash := hashToken(rawToken)

	session := &domainauth.RefreshSession{
		UserID:    userID,
		TokenHash: tokenHash,
		UserAgent: meta.UserAgent,
		IPAddress: meta.IPAddress,
		ExpiresAt: time.Now().Add(uc.cfg.RefreshTokenExpiry),
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}
	return rawToken, nil
}

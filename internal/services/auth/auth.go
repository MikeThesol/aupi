package auth

import (
	"auth-api/internal/services/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/MikeThesol/proto/proto/user/gen"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log        *slog.Logger
	userClient gen.UserClient // grpc-client for user service
	tokenTTL   time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *slog.Logger, userClient gen.UserClient, tokenTTL time.Duration) *Auth {
	return &Auth{
		userClient: userClient,
		log:        log,
		tokenTTL:   tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	userResponse, err := a.userClient.GetUserByEmail(ctx, &gen.GetUserByEmailRequest{Email: email})

	if err != nil {
		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userResponse.PassHash), []byte(password)); err != nil {
		log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(userResponse.Id, userResponse.IsAdmin, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err)
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	registerResponse, err := a.userClient.Register(ctx, &gen.RegisterRequest{
		Name:     name,
		Email:    email,
		PassHash: string(passHash),
	})

	if err != nil {
		log.Error("failed to register user", err)
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user registered successfully", slog.Int64("user_id", registerResponse.Id))

	return registerResponse.Id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	userResponse, err := a.userClient.GetUserByID(ctx, &gen.GetUserByIDRequest{Id: userID})
	if err != nil {
		log.Error("failed to get user", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", userResponse.IsAdmin))

	return userResponse.IsAdmin, nil
}

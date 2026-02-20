package tests

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/models"
	"sport-assistance/internal/services"
	"sport-assistance/pkg/commons"
	"sport-assistance/pkg/configs"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var errNotImplemented = errors.New("not implemented")

type mockRepository struct {
	createUserFn         func(ctx context.Context, user models.User) (uint64, error)
	getUserByIDFn        func(ctx context.Context, userID uint64) (models.User, error)
	getUserByEmailFn     func(ctx context.Context, email string) (models.User, error)
	userExistsByEmailFn  func(ctx context.Context, email string) (bool, error)
	rotateRefreshTokenFn func(ctx context.Context, userID uint64, oldRefreshToken, newRefreshToken string, newExpiresAt time.Time) error
	createRefreshTokenFn func(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error
	getRefreshTokenFn    func(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error)
	revokeRefreshTokenFn func(ctx context.Context, refreshToken string) error
}

func (m mockRepository) CreateUser(ctx context.Context, user models.User) (uint64, error) {
	if m.createUserFn == nil {
		return 0, errNotImplemented
	}
	return m.createUserFn(ctx, user)
}

func (m mockRepository) GetUserByID(ctx context.Context, userID uint64) (models.User, error) {
	if m.getUserByIDFn == nil {
		return models.User{}, errNotImplemented
	}
	return m.getUserByIDFn(ctx, userID)
}

func (m mockRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if m.getUserByEmailFn == nil {
		return models.User{}, errNotImplemented
	}
	return m.getUserByEmailFn(ctx, email)
}

func (m mockRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.userExistsByEmailFn == nil {
		return false, errNotImplemented
	}
	return m.userExistsByEmailFn(ctx, email)
}

func (m mockRepository) RotateRefreshToken(ctx context.Context, userID uint64, oldRefreshToken, newRefreshToken string, newExpiresAt time.Time) error {
	if m.rotateRefreshTokenFn == nil {
		return errNotImplemented
	}
	return m.rotateRefreshTokenFn(ctx, userID, oldRefreshToken, newRefreshToken, newExpiresAt)
}

func (m mockRepository) CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error {
	if m.createRefreshTokenFn == nil {
		return errNotImplemented
	}
	return m.createRefreshTokenFn(ctx, userID, refreshToken, expiresAt)
}

func (m mockRepository) GetRefreshToken(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error) {
	if m.getRefreshTokenFn == nil {
		return models.RefreshTokenResponse{}, errNotImplemented
	}
	return m.getRefreshTokenFn(ctx, refreshToken)
}

func (m mockRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if m.revokeRefreshTokenFn == nil {
		return errNotImplemented
	}
	return m.revokeRefreshTokenFn(ctx, refreshToken)
}

func testConfig() *configs.Config {
	return &configs.Config{
		DatabaseConfig: configs.DatabaseConfig{DBDateFormat: "02-01-2006"},
		SecurityConfig: configs.SecurityConfig{
			AccessTokenTTL:         time.Minute,
			AccessTokenSecret:      "access-secret",
			RefreshTokenTTL:        time.Hour,
			RefreshTokenSecret:     "refresh-secret",
			AccessTokenRedisPrefix: "auth:access_token:%d",
		},
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func unavailableRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:0",
		DialTimeout:  20 * time.Millisecond,
		ReadTimeout:  20 * time.Millisecond,
		WriteTimeout: 20 * time.Millisecond,
		PoolTimeout:  20 * time.Millisecond,
		MaxRetries:   0,
	})
}

func newService(repo services.IRepository) *services.Service {
	return services.NewService(repo, testLogger(), testConfig(), unavailableRedis())
}

func signRefreshToken(t *testing.T, cfg *configs.Config, userID uint64) string {
	t.Helper()

	now := time.Now()
	claims := models.CustomClaims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   commons.RefreshSubject,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.SecurityConfig.RefreshTokenSecret))
	if err != nil {
		t.Fatalf("failed to sign refresh token: %v", err)
	}
	return token
}

func TestUserExistsByEmail(t *testing.T) {
	expectedEmail := "user@example.com"
	service := newService(mockRepository{
		userExistsByEmailFn: func(_ context.Context, email string) (bool, error) {
			if email != expectedEmail {
				t.Fatalf("unexpected email: %s", email)
			}
			return true, nil
		},
	})

	exists, err := service.UserExistsByEmail(context.Background(), expectedEmail)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !exists {
		t.Fatalf("expected user to exist")
	}
}

func TestRegister_InvalidBirthDate(t *testing.T) {
	service := newService(mockRepository{})

	_, err := service.Register(context.Background(), requests.CreateUserRequest{
		Name:      "John",
		Surname:   "Doe",
		Gender:    "male",
		BirthDate: "bad-date",
		Email:     "john@example.com",
		Password:  "secret123",
	})

	if err == nil {
		t.Fatalf("expected birth date parse error")
	}
}

func TestRegister_CreateUserError(t *testing.T) {
	expectedErr := errors.New("db error")
	service := newService(mockRepository{
		createUserFn: func(_ context.Context, _ models.User) (uint64, error) {
			return 0, expectedErr
		},
	})

	_, err := service.Register(context.Background(), requests.CreateUserRequest{
		Name:      "John",
		Surname:   "Doe",
		Gender:    "male",
		BirthDate: "02-01-2000",
		Email:     "john@example.com",
		Password:  "secret123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	service := newService(mockRepository{
		getUserByEmailFn: func(_ context.Context, _ string) (models.User, error) {
			return models.User{}, errors.New("no rows")
		},
	})

	_, err := service.Login(context.Background(), requests.LoginRequest{Email: "user@example.com", Password: "pass"})
	if err == nil || err.Error() != "this user does not exist" {
		t.Fatalf("expected not exists error, got %v", err)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	service := newService(mockRepository{
		getUserByEmailFn: func(_ context.Context, _ string) (models.User, error) {
			return models.User{ID: 10, Email: "user@example.com", Password: string(hash)}, nil
		},
	})

	_, err = service.Login(context.Background(), requests.LoginRequest{Email: "user@example.com", Password: "wrong-password"})
	if err == nil || err.Error() != "invalid email or password" {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestLogout_EmptyRefreshToken(t *testing.T) {
	service := newService(mockRepository{})

	_, err := service.Logout(context.Background(), requests.LogoutRequest{UserID: 1})
	if err == nil || err.Error() != "refresh token is required" {
		t.Fatalf("expected refresh token required error, got %v", err)
	}
}

func TestLogout_TokenBelongsToAnotherUser(t *testing.T) {
	cfg := testConfig()
	service := services.NewService(mockRepository{}, testLogger(), cfg, unavailableRedis())
	refreshToken := signRefreshToken(t, cfg, 1)

	_, err := service.Logout(context.Background(), requests.LogoutRequest{UserID: 2, RefreshToken: refreshToken})
	if err == nil || err.Error() != "refresh token does not belong to this user" {
		t.Fatalf("expected user mismatch error, got %v", err)
	}
}

func TestCreateTokens_RedisUnavailable(t *testing.T) {
	service := newService(mockRepository{})

	_, _, err := service.CreateTokens(context.Background(), 1, "user@example.com")
	if err == nil {
		t.Fatalf("expected redis error")
	}
}

func TestRefreshTokens_InvalidRefreshToken(t *testing.T) {
	service := newService(mockRepository{})

	_, err := service.RefreshTokens(context.Background(), requests.RefreshTokensRequest{RefreshToken: "not-a-token"})
	if err == nil || err.Error() != "invalid refresh token" {
		t.Fatalf("expected invalid refresh token error, got %v", err)
	}
}

func TestRefreshTokens_RevokedToken(t *testing.T) {
	cfg := testConfig()
	token := signRefreshToken(t, cfg, 15)
	revokedAt := time.Now()

	service := services.NewService(mockRepository{
		getRefreshTokenFn: func(_ context.Context, _ string) (models.RefreshTokenResponse, error) {
			return models.RefreshTokenResponse{
				UserID:    15,
				ExpiresAt: time.Now().Add(time.Hour),
				RevokedAt: &revokedAt,
			}, nil
		},
	}, testLogger(), cfg, unavailableRedis())

	_, err := service.RefreshTokens(context.Background(), requests.RefreshTokensRequest{RefreshToken: token})
	if err == nil || err.Error() != "refresh token is revoked" {
		t.Fatalf("expected revoked token error, got %v", err)
	}
}

func TestIsTokenExpired(t *testing.T) {
	service := newService(mockRepository{})

	if !service.IsTokenExpired(time.Now().Add(-time.Second)) {
		t.Fatalf("expected past token to be expired")
	}
	if service.IsTokenExpired(time.Now().Add(time.Minute)) {
		t.Fatalf("expected future token to be active")
	}
}

package handlers

import (
	"context"
	"log/slog"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/pkg/configs"

	"github.com/gin-gonic/gin"
)

type IService interface {

	// jwt
	Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error)
	Login(ctx context.Context, req requests.LoginRequest) (responses.JWTResponse, error)
	CreateTokens(ctx context.Context, userID uint64, email string) (string, string, error)
	RefreshTokens(ctx context.Context, request requests.RefreshTokensRequest) (responses.JWTResponse, error)
	Logout(ctx context.Context, request requests.LogoutRequest) (responses.EmptyResponse, error)

	//OTP
	SendOTP(ctx context.Context, identifier string) (responses.SendOTPResponse, error)
	ConfirmOTP(ctx context.Context, identifier, otp string) (responses.ConfirmOTPResponse, error)
}
type IMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
	CORSMiddleware() gin.HandlerFunc
	RequirePermissions(permissions ...string) gin.HandlerFunc
}

type Handler struct {
	service     IService
	logger      *slog.Logger
	middlewares IMiddleware
	cfg         *configs.Config
}

func NewHandler(service IService, log *slog.Logger, middlewares IMiddleware, cfg *configs.Config) *Handler {
	return &Handler{
		service:     service,
		logger:      log,
		middlewares: middlewares,
		cfg:         cfg,
	}
}

func (h *Handler) InitHandler() *gin.Engine {
	router := gin.New()
	router.Use(h.middlewares.CORSMiddleware(), gin.RecoveryWithWriter(gin.DefaultWriter))

	ping := router.Group("/")
	{
		ping.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	public := router.Group("/api/v1/auth")
	{
		public.POST("/registration", h.Register)
		public.POST("/login", h.Login)
		public.POST("/refresh", h.RefreshTokens)
		public.POST("/logout", h.Logout)
		public.POST("/otp/send", h.SendOTP)
		public.POST("/otp/confirm", h.ConfirmOTP)
	}

	private := router.Group("/api/v1")
	private.Use(h.middlewares.AuthMiddleware())
	private.Use(h.middlewares.CORSMiddleware())

	profile := private.Group("/profile")
	profile.Use(h.middlewares.RequirePermissions("profile.view.own"))
	{
		profile.GET("/me", func(c *gin.Context) {})
	}

	match := private.Group("/match")
	match.Use(
		h.middlewares.RequirePermissions(
			"match.confirm.participation",
			"match.create",
			"match.enter.result",
			"match.invite.users"),
	)
	{
	}

	return router
}

package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/pkg/configs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type IService interface {
	Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error)
	Login(ctx context.Context, req requests.LoginRequest) (responses.JWTResponse, error)
	CreateTokens(ctx context.Context, userID uint64, email string) (string, string, error)
	RefreshTokens(ctx context.Context, request requests.RefreshTokensRequest) (responses.JWTResponse, error)
}
type IMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
	CORSMiddleware() gin.HandlerFunc
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

	if h.cfg.SwaggerConfig.SwaggerEnabled {
		router.StaticFile("/swagger.yaml", "./swagger.yaml")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("/openapi.yaml"),
		))
	}

	public := router.Group("/api/v1/auth")
	{
		public.POST("/registration", h.Register)
		public.POST("/login", h.Login)
		public.POST("/refresh", h.RefreshTokens)
		public.POST("/logout")
	}

	private := router.Group("/api/v1")
	private.Use(h.middlewares.AuthMiddleware())
	private.Use(h.middlewares.CORSMiddleware())
	{
		private.GET("/lolo", func(c *gin.Context) {
			fmt.Println("lolo")
		})
	}
	return router
}

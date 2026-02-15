package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"

	"github.com/gin-gonic/gin"
)

type IService interface {
	Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error)
	CreateTokens(ctx context.Context, userID int64, email string) (string, string, error)
}
type IMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
	CORSMiddleware() gin.HandlerFunc
}

type Handler struct {
	service     IService
	logger      *slog.Logger
	middlewares IMiddleware
}

func NewHandler(service IService, log *slog.Logger, middlewares IMiddleware) *Handler {
	return &Handler{
		service:     service,
		logger:      log,
		middlewares: middlewares,
	}
}

func (h *Handler) InitHandler() *gin.Engine {
	router := gin.New()

	router.Use(h.middlewares.CORSMiddleware(), gin.RecoveryWithWriter(gin.DefaultWriter))

	public := router.Group("/api/v1")
	{
		public.POST("/registration", h.Register)
		public.POST("/login")
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

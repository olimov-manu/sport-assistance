package handlers

import (
	"sport-assistance/internal/handlers/requests"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind error", "err", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(ctx, req)
	if err != nil {
		h.logger.Error("register failed", "err", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind error", "err", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jwts, err := h.service.Login(ctx, req)
	if err != nil {
		h.logger.Error("login failed", "err", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, jwts)
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.RefreshTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind error", "err", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, err := h.service.RefreshTokens(ctx, req)
	if err != nil {
		h.logger.Error("refresh access token failed", "err", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, newAccessToken)
}

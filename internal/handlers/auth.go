package handlers

import (
	"errors"
	"net/http"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/pkg/myerrors"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleError(c *gin.Context, err error) {
	var appErr myerrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(http.StatusBadRequest, appErr.ToResponse())
		return
	}
	c.JSON(http.StatusInternalServerError, myerrors.Response{
		Message: "Internal server error",
		Error:   err.Error(),
	})
}

func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Create user request bind error ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.service.Register(ctx, req)
	if err != nil {
		h.logger.Error("Registration failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind login request error: ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	jwts, err := h.service.Login(ctx, req)
	if err != nil {
		h.logger.Error("Login failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, jwts)
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.RefreshTokensRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind tokens error: ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	newAccessToken, err := h.service.RefreshTokens(ctx, req)
	if err != nil {
		h.logger.Error("Refresh tokens failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, newAccessToken)
}

func (h *Handler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind logout request error: ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if _, err := h.service.Logout(ctx, req); err != nil {
		h.logger.Error("Logout failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

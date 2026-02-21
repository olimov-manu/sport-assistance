package handlers

import (
	"net/http"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/pkg/myerrors"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendOTP(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind send otp request error: ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.service.SendOTP(ctx, req.Identifier)
	if err != nil {
		h.logger.Error("Send otp failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ConfirmOTP(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.ConfirmOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind confirm otp request error: ", "err", err)
		c.JSON(http.StatusBadRequest, myerrors.Response{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.service.ConfirmOTP(ctx, req.Identifier, req.OTP)
	if err != nil {
		h.logger.Error("Confirm otp failed: ", "err", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

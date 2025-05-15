package http

import (
	"net/http"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc    port.UserService
	logger *zap.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
	}
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.logger.Warn("Failed to bind register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	createdUser, err := uh.svc.Register(c.Request.Context(), &req)
	if err != nil {
		uh.logger.Error("Registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, util.APIResponse("Registration failed", http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Registration successful", http.StatusOK, "success", createdUser))
}

func (uh *UserHandler) GetProfile(c *gin.Context) {
	payload := util.GetAuthPayload(c, consts.AuthorizationKey)
	if payload == nil {
		uh.logger.Error("No authorization payload found")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	id := payload.UserID
	user, err := uh.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		uh.logger.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusNotFound, util.APIResponse("User not found", http.StatusNotFound, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("User profile fetched", http.StatusOK, "success", user))
}

func (uh *UserHandler) UpdateProfile(c *gin.Context) {
	payload := util.GetAuthPayload(c, consts.AuthorizationKey)
	if payload == nil {
		uh.logger.Error("No authorization payload found")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.logger.Warn("Failed to bind update profile request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	updatedUser, err := uh.svc.UpdateProfile(c, req)
	if err != nil {
		uh.logger.Error("Failed to update user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, util.APIResponse("Failed to update user profile", http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("User profile updated", http.StatusOK, "success", updatedUser))
}

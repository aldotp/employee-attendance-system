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

func (h *UserHandler) ListUser(c *gin.Context) {
	var req dto.ListUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.svc.ListUsers(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("List Users Success", http.StatusOK, "success", users))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.svc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Get User Success", http.StatusOK, "success", user))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		if err == consts.ErrEmailAlreadyExist {
			c.JSON(http.StatusBadRequest, util.APIResponse("Email already exist", http.StatusBadRequest, "error", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusCreated, util.APIResponse("Create User Success", http.StatusCreated, "success", user))
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.UpdateUserByID(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Update User Success", http.StatusOK, "success", user))
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Delete User Success", http.StatusOK, "success", nil))
}

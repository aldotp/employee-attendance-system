package http

import (
	"net/http"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	svc port.NotificationService
}

func NewNotificationHandler(svc port.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		svc: svc,
	}
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req domain.Notification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	id, err := h.svc.CreateNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse("Notification created", http.StatusCreated, "success", gin.H{"id": id}))
}

func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	id := c.Param("id")
	notification, err := h.svc.GetNotificationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", notification))
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userSession := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSession == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input dto.ListNotificationRequest
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	notifications, err := h.svc.ListNotifications(c.Request.Context(), userSession.UserID, input.Skip, input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", notifications))
}

func (h *NotificationHandler) UpdateNotificationStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	err := h.svc.UpdateNotificationStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.Status(http.StatusOK)
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.Status(http.StatusOK)
}

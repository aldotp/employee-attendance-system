package http

import (
	"net/http"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	svc port.LeaveService
}

func NewLeaveHandler(svc port.LeaveService) *LeaveHandler {
	return &LeaveHandler{svc: svc}
}

func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	leaves, err := h.svc.ListLeaves(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("success", http.StatusOK, "success", leaves))
}

func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	userSession := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSession == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.EmployeeID = userSession.EmployeeID
	leave, err := h.svc.SubmitLeaveRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leave)
}

func (h *LeaveHandler) GetLeave(c *gin.Context) {
	id := c.Param("id")
	leave, err := h.svc.GetLeaveByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leave)
}

func (h *LeaveHandler) UpdateLeave(c *gin.Context) {
	id := c.Param("id")
	var req dto.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	leave, err := h.svc.UpdateLeave(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leave)
}

func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteLeave(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *LeaveHandler) GetLeaveBalance(c *gin.Context) {
	var req dto.LeaveBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.svc.GetLeaveBalance(c.Request.Context(), req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	leaveID := c.Param("id")

	err := h.svc.ApproveLeave(c, leaveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Approve Leave Request Success", http.StatusOK, "success", nil))
}

func (h *LeaveHandler) RejectLeave(c *gin.Context) {
	leaveID := c.Param("id")

	var req dto.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.RejectLeave(c, leaveID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Reject Leave Request Success", http.StatusOK, "success", nil))
}

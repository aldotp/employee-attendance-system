package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/service"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(svc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

func (h *AttendanceHandler) ListAttendance(c *gin.Context) {
	page := uint64(1)
	limit := uint64(10)
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	attendances, err := h.svc.ListAttendances(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("List Attendances", http.StatusOK, "success", attendances))
}

func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	payload := util.GetAuthPayload(c, consts.AuthorizationKey)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	var req dto.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.OpenAttendance(c, req, payload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	id := c.Param("id")
	attendance, err := h.svc.GetAttendanceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	id := c.Param("id")
	var req dto.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	attendance, err := h.svc.UpdateAttendance(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) DeleteAttendance(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteAttendance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AttendanceHandler) GetAttendanceHistory(c *gin.Context) {
	employeeID := c.Query("employee_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	history, err := h.svc.GetAttendanceHistory(c, employeeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *AttendanceHandler) GetUsersAttendanceStatus(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	statusMap, err := h.svc.GetUsersAttendanceStatus(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":  date,
		"users": statusMap,
	})
}

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

type ScheduleHandler struct {
	scheduleService port.ScheduleService
}

func NewScheduleHandler(scheduleService port.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) ListSchedules(c *gin.Context) {
	schedules, err := h.scheduleService.ListSchedules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success List Schedule", http.StatusOK, "success", schedules))
}
func (h *ScheduleHandler) GetWorkRotation(c *gin.Context) {
	employeeID := c.Param("employeeId")
	rotation, err := h.scheduleService.GetWorkRotation(c.Request.Context(), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success Get Work Rotation", http.StatusOK, "success", rotation))
}

func (h *ScheduleHandler) GetWorkCalendar(c *gin.Context) {
	userSession := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSession == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input dto.GetWorkCalendar
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	calendar, err := h.scheduleService.GetWorkCalendar(c, userSession.UserID, input.Year, input.Month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, calendar)
}

func (h *ScheduleHandler) RequestScheduleSwap(c *gin.Context) {
	userSession := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSession == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.ScheduleSwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.scheduleService.RequestScheduleSwap(c, userSession.UserID, req.ScheduleID1, req.ScheduleID2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req domain.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.scheduleService.CreateSchedule(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, util.APIResponse("Schedule created", http.StatusCreated, "success", schedule))
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	schedule, err := h.scheduleService.GetSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", schedule))
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	var req domain.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.scheduleService.UpdateSchedule(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", schedule))
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	err := h.scheduleService.DeleteSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Delete Schedule Success", http.StatusOK, "success", nil))
}

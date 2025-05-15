package http

import (
	"net/http"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	svc port.MonitoringService
}

func NewMonitoringHandler(svc port.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{
		svc: svc,
	}
}

func (h *MonitoringHandler) GetReports(c *gin.Context) {
	reports, err := h.svc.GetReports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", reports))
}

func (h *MonitoringHandler) GetSummary(c *gin.Context) {
	summary, err := h.svc.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", summary))
}

func (h *MonitoringHandler) GetDashboardAnalytics(c *gin.Context) {
	analytics, err := h.svc.GetDashboardAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", analytics))
}

func (h *MonitoringHandler) GenerateAttendanceReport(c *gin.Context) {
	report, err := h.svc.GenerateAttendanceReport(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", report))
}

func (h *MonitoringHandler) DetectAnomalies(c *gin.Context) {
	anomalies, err := h.svc.DetectAnomalies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", anomalies))
}

func (h *MonitoringHandler) ExportData(c *gin.Context) {
	var req domain.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	data, err := h.svc.ExportData(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", data))
}

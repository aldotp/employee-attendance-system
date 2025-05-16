package http

import (
	"net/http"
	"strconv"
	"time"

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
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	summary, err := h.svc.GetSummary(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("Success", http.StatusOK, "success", summary))
}

func (h *MonitoringHandler) GetDashboardAnalytics(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	analytics, err := h.svc.GetDashboardAnalytics(c.Request.Context(), date)
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
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	_, StartDateErr := time.Parse("2006-01-02", req.StartDate)
	_, EndDateErr := time.Parse("2006-01-02", req.EndDate)

	if StartDateErr != nil || EndDateErr != nil {
		errorMessage := gin.H{"errors": "please input correct date"}
		response := util.APIResponse("export data failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	exportResponse, err := h.svc.ExportData(c.Request.Context(), req)
	if err != nil {
		response := util.APIResponse("failed get data", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	downloadName := exportResponse.FileName
	bytes := exportResponse.FileContent

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", strconv.Itoa(bytes.Len()))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=\""+downloadName+"\"")
	c.Data(http.StatusOK, "application/octet-stream", bytes.Bytes())
}

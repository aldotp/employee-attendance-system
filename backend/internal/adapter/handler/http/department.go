package http

import (
	"net/http"
	"strconv"

	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentService port.DepartmentService
}

func NewDepartmentHandler(departmentService port.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	var limit, skip uint64

	skipQuery := c.Query("skip")
	if skipQuery == "" {
		skip = 0
	} else {
		if val, err := strconv.ParseUint(skipQuery, 10, 64); err == nil {
			skip = val
		} else {
			skip = 0
		}
	}

	limitQuery := c.Query("limit")
	if limitQuery == "" {
		limit = 10
	} else {
		if val, err := strconv.ParseUint(limitQuery, 10, 64); err == nil {
			limit = val
		} else {
			limit = 10
		}
	}

	departments, err := h.departmentService.ListDepartments(c.Request.Context(), skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse("success", http.StatusOK, "list of departments", departments))
}

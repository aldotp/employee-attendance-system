package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseIntParam safely extracts and converts a path parameter to an integer.
// Returns 0 if the parameter is missing or invalid.
func ParseIntParam(c *gin.Context, param string) int {
	valueStr := c.Param(param)
	if valueStr == "" {
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}

	return value
}

package util

import (
	"context"
	"strconv"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

func StringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)
	return num, err
}

func GetAuthPayload(ctx context.Context, key string) *domain.TokenPayload {
	val := ctx.Value(key)
	if payload, ok := val.(*domain.TokenPayload); ok {
		return payload
	}
	return nil
}

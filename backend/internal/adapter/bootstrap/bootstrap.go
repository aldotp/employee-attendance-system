package bootstrap

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/config"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/gcs"
	"github.com/aldotp/employee-attendance-system/pkg/minio"

	"go.uber.org/zap"
)

type Bootstrap struct {
	ctx        context.Context
	Log        *zap.Logger
	Config     *config.Config
	PostgresDB *postgres.DB
	GCS        *gcs.GCS
	minio      *minio.MinioClient

	AttendanceRepo   port.AttendanceRepository
	DepartmentRepo   port.DepartmentRepository
	DeviceLogRepo    port.DeviceLogRepository
	DeviceRepo       port.DeviceRepository
	EmployeeRepo     port.EmployeeRepository
	LeaveRequestRepo port.LeaveRequestRepository
	NotificationRepo port.NotificationRepository
	UserRepo         port.UserRepository
	WorkLocationRepo port.WorkLocationRepository
	ScheduleRepo     port.ScheduleRepository
	MonitoringRepo   port.MonitoringRepository

	Token port.TokenInterface
	Cache port.CacheInterface
}

func NewBootstrap(ctx context.Context) *Bootstrap {
	return &Bootstrap{
		ctx: ctx,
	}
}

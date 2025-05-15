package bootstrap

import (
	"log"
	"log/slog"
	"os"

	"github.com/aldotp/employee-attendance-system/internal/adapter/auth/jwt"
	"github.com/aldotp/employee-attendance-system/internal/adapter/config"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	postgresRepo "github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres/repository"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/redis"
	"github.com/aldotp/employee-attendance-system/pkg/gcs"
	"github.com/aldotp/employee-attendance-system/pkg/logger"
	"github.com/aldotp/employee-attendance-system/pkg/minio"
)

func (b *Bootstrap) setConfig() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	b.Config = config
}

func (b *Bootstrap) setLogger() {
	logger, err := logger.InitLogger(config.AppEnv())
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	b.Log = logger
}

func (b *Bootstrap) setPostgresDB() {
	db, err := postgres.New(b.ctx, &config.DB{
		Connection: config.DBConnection(),
		User:       config.DBUser(),
		Password:   config.DBPassword(),
		Host:       config.DBHost(),
		Port:       config.DBPort(),
		Name:       config.DBName(),
	})
	if err != nil {
		panic(err)
	}

	b.PostgresDB = db
}

func (b *Bootstrap) setJWTToken() {
	token, err := jwt.New()
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	b.Token = token
}

func (b *Bootstrap) setCache() {
	cache, err := redis.New(b.ctx, &config.Redis{
		Addr:     config.RedisAddr(),
		Password: config.RedisPassword(),
	})
	if err != nil {
		panic(err)
	}

	b.Cache = cache
}

func (b *Bootstrap) setRestApiRepository() {
	b.UserRepo = postgresRepo.NewUserRepository(b.PostgresDB)
	b.AttendanceRepo = postgresRepo.NewAttendanceRepository(b.PostgresDB)
	b.DepartmentRepo = postgresRepo.NewDepartmentRepository(b.PostgresDB)
	b.DeviceLogRepo = postgresRepo.NewDeviceLogRepository(b.PostgresDB)
	b.DeviceRepo = postgresRepo.NewDeviceRepository(b.PostgresDB)
	b.EmployeeRepo = postgresRepo.NewEmployeeRepository(b.PostgresDB)
	b.LeaveRequestRepo = postgresRepo.NewLeaveRequestRepository(b.PostgresDB)
	b.NotificationRepo = postgresRepo.NewNotificationRepository(b.PostgresDB)
	b.WorkLocationRepo = postgresRepo.NewWorkLocationRepository(b.PostgresDB)
	b.ScheduleRepo = postgresRepo.NewScheduleRepository(b.PostgresDB)
	b.MonitoringRepo = postgresRepo.NewMonitoringRepository(b.PostgresDB)
}

func (b *Bootstrap) setGCS() {
	gcs, err := gcs.NewCGS("../../../creds/gcs.json", "employee-attendance-system")
	if err != nil {
		log.Fatalf("error gcs %v", err.Error())
	}

	b.GCS = gcs

}

func (b *Bootstrap) SetMinio() {
	minio, err := minio.NewMinioClient(config.MinioEndpoint(), config.MinioAccessKey(), config.MinioSecretKey(), config.MinioBucketName(), config.MinioUseSSL())
	if err != nil {
		log.Fatalf("error minio %v", err.Error())
	}

	b.minio = minio
}

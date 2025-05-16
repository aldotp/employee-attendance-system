package http

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aldotp/employee-attendance-system/internal/adapter/bootstrap"
	"github.com/aldotp/employee-attendance-system/internal/adapter/handler/http"
	"github.com/aldotp/employee-attendance-system/internal/adapter/router"
	"github.com/aldotp/employee-attendance-system/internal/core/service"
)

func RunHTTPServer(ctx context.Context) {
	f := bootstrap.NewBootstrap(ctx).BuildRestBootstrap()

	// Services
	userService := service.NewUserService(f.UserRepo, f.EmployeeRepo, f.DepartmentRepo, f.Cache, f.Token, f.Log)
	authService := service.NewAuthService(f.UserRepo, f.EmployeeRepo, f.Token, f.Log)
	attendanceService := service.NewAttendanceService(f.AttendanceRepo)
	leaveService := service.NewLeaveService(f.LeaveRequestRepo, f.NotificationRepo)
	scheduleService := service.NewScheduleService(f.ScheduleRepo)
	monitoringService := service.NewMonitoringService(f.MonitoringRepo, f.UserRepo, f.AttendanceRepo)
	notificationService := service.NewNotificationService(f.NotificationRepo, f.UserRepo)

	// Handlers
	userHandler := http.NewUserHandler(userService, f.Log)
	authHandler := http.NewAuthHandler(authService, f.Log)
	attendanceHandler := http.NewAttendanceHandler(attendanceService)
	leaveHandler := http.NewLeaveHandler(leaveService)
	scheduleHandler := http.NewScheduleHandler(scheduleService)
	monitoringHandler := http.NewMonitoringHandler(monitoringService)
	notificationHandler := http.NewNotificationHandler(notificationService)
	deparmentHandler := http.NewDepartmentHandler(f.DepartmentRepo)

	// HTTP server
	routes, err := router.NewRouter(
		f.Token,
		authHandler,
		userHandler,
		attendanceHandler,
		leaveHandler,
		scheduleHandler,
		monitoringHandler,
		notificationHandler,
		deparmentHandler,
	)
	if err != nil {
		slog.Error("Error creating router", "error", err)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", f.Config.HTTP.URL, f.Config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = routes.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}

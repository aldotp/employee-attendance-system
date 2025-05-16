package router

import (
	"github.com/aldotp/employee-attendance-system/internal/adapter/config"
	"github.com/aldotp/employee-attendance-system/internal/adapter/handler/http"
	"github.com/aldotp/employee-attendance-system/internal/adapter/middleware"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	token port.TokenInterface,
	authHandler *http.AuthHandler,
	userHandler *http.UserHandler,
	attendanceHandler *http.AttendanceHandler,
	leaveHandler *http.LeaveHandler,
	scheduleHandler *http.ScheduleHandler,
	monitoringHandler *http.MonitoringHandler,
	notificationHandler *http.NotificationHandler,
) (*Router, error) {

	// Set Gin mode
	if config.AppEnv() == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Custom validator
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
			return nil, err
		}
	}

	util.Index(router, config.AppVersion(), config.AppName())
	util.Metrics(router)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{

		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		user := v1.Group("/user").Use(middleware.AuthMiddleware(token))
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		admin := v1.Group("/admin").Use(middleware.AuthMiddleware(token), middleware.VerifyRole(domain.Admin, domain.HR))
		{
			admin.GET("/users", userHandler.ListUser)
			admin.GET("/users/:id", userHandler.GetUserByID)
			admin.POST("/users", userHandler.CreateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUserByID)
			admin.PUT("/users/:id", userHandler.UpdateUserByID)
		}

		notification := v1.Group("/notification").Use(middleware.AuthMiddleware(token))
		{
			notification.GET("", notificationHandler.ListNotifications)
			notification.GET("/:id", notificationHandler.GetNotificationByID)
			notification.PUT("/:id", notificationHandler.UpdateNotificationStatus)
			notification.DELETE("/:id", notificationHandler.DeleteNotification)
			notification.POST("", notificationHandler.CreateNotification)
		}

		attendance := v1.Group("/attendance")
		{
			att := attendance.Use(middleware.AuthMiddleware(token))
			att.GET("", attendanceHandler.ListAttendance)
			att.POST("", attendanceHandler.CreateAttendance)
			att.GET("/:id", attendanceHandler.GetAttendance)
			att.PUT("/:id", attendanceHandler.UpdateAttendance)
			att.DELETE("/:id", attendanceHandler.DeleteAttendance)
			att.GET("/status", attendanceHandler.GetUsersAttendanceStatus)
		}

		leave := v1.Group("/leave")
		{
			leaveUser := leave.Group("").Use(middleware.AuthMiddleware(token))
			leaveUser.GET("", leaveHandler.ListLeaves)
			leaveUser.POST("", leaveHandler.CreateLeave)
			leaveUser.GET("/:id", leaveHandler.GetLeave)
			leaveUser.PUT("/:id", leaveHandler.UpdateLeave)
			leaveUser.DELETE("/:id", leaveHandler.DeleteLeave)

			leaveAdmin := leave.Group("/admin").Use(middleware.AuthMiddleware(token), middleware.VerifyRole(domain.Admin, domain.HR))
			leaveAdmin.GET("/balance", leaveHandler.GetLeaveBalance)
			leaveAdmin.POST("/approve/:id", leaveHandler.ApproveLeave)
			leaveAdmin.POST("/reject/:id", leaveHandler.RejectLeave)
		}

		schedule := v1.Group("/schedule").Use(middleware.AuthMiddleware(token))
		{
			schedule.GET("", scheduleHandler.ListSchedules)
			schedule.POST("", scheduleHandler.CreateSchedule)
			schedule.GET("/:id", scheduleHandler.GetSchedule)
			schedule.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedule.DELETE("/:id", scheduleHandler.DeleteSchedule)
			schedule.GET("/rotation", scheduleHandler.GetWorkRotation)
			schedule.GET("/calendar", scheduleHandler.GetWorkCalendar)
			schedule.POST("/swap", scheduleHandler.RequestScheduleSwap)
		}

		monitoring := v1.Group("/monitoring").Use(middleware.AuthMiddleware(token))
		{
			monitoring.GET("/reports", monitoringHandler.GetReports)
			monitoring.GET("/summary", monitoringHandler.GetSummary)
			monitoring.GET("/dashboard", monitoringHandler.GetDashboardAnalytics)
			monitoring.GET("/attendance-report", monitoringHandler.GenerateAttendanceReport)
			monitoring.GET("/export", monitoringHandler.ExportData)
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}

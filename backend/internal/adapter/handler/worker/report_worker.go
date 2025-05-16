package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/bootstrap"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/internal/core/service"
	"github.com/google/uuid"
	"github.com/mileusna/crontab"
)

type ReportWorker struct {
	attendanceService port.AttendanceService
	monitoringService port.MonitoringService
	monitoringRepo    port.MonitoringRepository
	ctab              *crontab.Crontab
}

func NewReportWorker(b *bootstrap.Bootstrap) *ReportWorker {
	return &ReportWorker{
		attendanceService: service.NewAttendanceService(b.AttendanceRepo),
		monitoringService: service.NewMonitoringService(b.MonitoringRepo, b.UserRepo, b.AttendanceRepo),
		monitoringRepo:    b.MonitoringRepo,
		ctab:              crontab.New(),
	}
}

func (w *ReportWorker) Run() {
	err := w.ctab.AddJob("* * * * *", w.GenerateSummaryReport)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Scheduler Running: Daily Attendance Report Generation")
	}

	err = w.ctab.AddJob("0 0 31 * *", w.GenerateAttendancyReport)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Scheduler Running: Daily Attendance Report Generation")
	}
}

func (w *ReportWorker) GenerateSummaryReport() {
	ctx := context.Background()
	tNow := time.Now()
	log.Println("Generating Summary report...")

	summary, err := w.monitoringService.GetSummary(ctx, tNow.Format("2006-01-02"))
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	data, err := json.Marshal(summary)
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	err = w.monitoringRepo.InsertReport(ctx, &domain.MonitoringReport{
		ID:          uuid.NewString(),
		ReportType:  "summary",
		Data:        string(data),
		GeneratedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	log.Println("Summary report generated successfully", time.Now())
}

func (w *ReportWorker) GenerateAttendancyReport() {
	ctx := context.Background()
	log.Println("Generating Attendence report...")

	attendanceData, err := w.monitoringService.GenerateAttendanceReport(ctx)
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	data, err := json.Marshal(attendanceData)
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	err = w.monitoringRepo.InsertReport(ctx, &domain.MonitoringReport{
		ID:          uuid.NewString(),
		ReportType:  "attendance",
		Data:        string(data),
		GeneratedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		log.Println("Failed to generate attendance report:", err)
		return
	}

	log.Println("Attendance report generated successfully", time.Now())
}

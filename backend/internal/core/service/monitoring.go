package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/xuri/excelize/v2"
)

type MonitoringService struct {
	repo           port.MonitoringRepository
	userRepo       port.UserRepository
	attendanceRepo port.AttendanceRepository
}

func NewMonitoringService(repo port.MonitoringRepository, userRepo port.UserRepository, attendanceRepo port.AttendanceRepository) *MonitoringService {
	return &MonitoringService{
		repo:           repo,
		userRepo:       userRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (ms *MonitoringService) GetReports(ctx context.Context) ([]domain.MonitoringReport, error) {
	return ms.repo.GetReports(ctx)
}

func (ms *MonitoringService) GetSummary(ctx context.Context, date string) (*domain.MonitoringSummary, error) {
	return ms.repo.GetSummary(ctx, date)
}

func (ms *MonitoringService) GetDashboardAnalytics(ctx context.Context, date string) (*domain.DashboardAnalytics, error) {
	summary, err := ms.repo.GetSummary(ctx, date)
	if err != nil {
		return nil, err
	}

	var startOfWeek time.Time
	if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		weekday := int(parsedDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startOfWeek = parsedDate.AddDate(0, 0, -(weekday - 1))
	} else {
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startOfWeek = now.AddDate(0, 0, -(weekday - 1))
	}

	weeklyAttendance := make([]int, 7)
	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		dayStr := day.Format("2006-01-02")
		summaryDay, err := ms.repo.GetSummary(ctx, dayStr)
		if err != nil {
			weeklyAttendance[i] = 0
			continue
		}
		weeklyAttendance[i] = summaryDay.TotalCheckin + summaryDay.TotalCheckOut
	}

	analytics := &domain.DashboardAnalytics{
		WeeklyAttendance: weeklyAttendance,
		LeaveDistribution: map[string]int{
			"pending":  summary.PendingLeaves,
			"approved": summary.ApprovedLeaves,
			"rejected": summary.RejectedLeaves,
		},
		GeneratedAt: time.Now(),
	}

	return analytics, nil
}

func (ms *MonitoringService) GenerateAttendanceReport(ctx context.Context) ([]domain.AttendanceReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	users, err := ms.userRepo.FindAllWithDetails(ctx)
	if err != nil {
		return nil, err
	}

	var (
		reportsMu sync.Mutex
		reports   []domain.AttendanceReport
		wg        sync.WaitGroup
		errChan   = make(chan error, len(users))
	)

	daysInMonth := int(endOfMonth.Day())

	for _, user := range users {
		wg.Add(1)
		go func(user domain.UserWithEmployee) {
			defer wg.Done()

			attendances, err := ms.attendanceRepo.GetAttendanceHistory(ctx, user.ID, startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"))
			if err != nil {
				errChan <- err
				return
			}

			attendanceMap := make(map[string]domain.AttendanceStatus)
			for _, att := range attendances {
				dayStr := att.Time.Format("2006-01-02")
				attendanceMap[dayStr] = att.Status
			}

			var lateCount, absentCount int
			dailyStatus := make(map[string]domain.AttendanceStatus)

			for i := 0; i < daysInMonth; i++ {
				day := startOfMonth.AddDate(0, 0, i)
				dayStr := day.Format("2006-01-02")
				status, ok := attendanceMap[dayStr]
				if !ok {
					status = domain.AttendanceStatusAbsent
					absentCount++
				} else {
					if status == domain.AttendanceStatusLate {
						lateCount++
					}
				}
				dailyStatus[dayStr] = status
			}

			report := domain.AttendanceReport{
				Name:        user.Name,
				UserID:      user.ID,
				LateCount:   lateCount,
				AbsentCount: absentCount,
				PeriodStart: startOfMonth,
				PeriodEnd:   endOfMonth,
				DailyStatus: dailyStatus,
			}

			reportsMu.Lock()
			reports = append(reports, report)
			reportsMu.Unlock()
		}(user)
	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return reports, nil
}

func (ms *MonitoringService) DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error) {
	return []domain.Anomaly{}, nil
}

func (ms *MonitoringService) ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error) {
	data, err := ms.repo.GetReportByDate(ctx, req)
	if err != nil {
		return nil, err
	}

	excelHandle := excelize.NewFile()
	defer func() {
		if err = excelHandle.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	header := []string{
		"Date",
		"Report Type",
		"Total Users",
		"Active Users",
		"Total Check_In",
		"Total Check_Out",
		"Pending Leaves",
		"Approved Leaves",
		"Rejected Leaves",
	}

	sheet := "Report"
	index, err := excelHandle.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	excelHandle.SetActiveSheet(index)

	// Write header
	for i, v := range header {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, err
		}
		excelHandle.SetCellValue(sheet, cell, v)
	}

	// Write data
	for i, report := range data {
		var summary domain.MonitoringSummary
		if err := json.Unmarshal([]byte(report.Data), &summary); err != nil {
			return nil, err
		}

		row := i + 2
		excelHandle.SetCellValue(sheet, fmt.Sprintf("A%d", row), report.GeneratedAt.Format("2006-01-02"))
		excelHandle.SetCellValue(sheet, fmt.Sprintf("B%d", row), report.ReportType)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("C%d", row), summary.TotalUsers)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("D%d", row), summary.ActiveUsers)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("E%d", row), summary.TotalCheckin)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("F%d", row), summary.TotalCheckOut)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("G%d", row), summary.PendingLeaves)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("H%d", row), summary.ApprovedLeaves)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("I%d", row), summary.RejectedLeaves)
	}

	var b bytes.Buffer
	if err := excelHandle.Write(&b); err != nil {
		return nil, err
	}

	return &domain.ExportResponse{
		FileContent: b,
		FileName:    fmt.Sprintf("report_%s.xlsx", time.Now().Format("20060102_150405")),
	}, nil
}

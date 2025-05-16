package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/xuri/excelize/v2"
)

type MonitoringService struct {
	repo port.MonitoringRepository
}

func NewMonitoringService(repo port.MonitoringRepository) *MonitoringService {
	return &MonitoringService{
		repo: repo,
	}
}

func (ms *MonitoringService) GetReports(ctx context.Context) ([]domain.MonitoringReport, error) {
	return ms.repo.GetReports(ctx)
}

func (ms *MonitoringService) GetSummary(ctx context.Context, date string) (*domain.MonitoringSummary, error) {
	return ms.repo.GetSummary(ctx, date)
}

func (ms *MonitoringService) GetDashboardAnalytics(ctx context.Context, date string) (*domain.DashboardAnalytics, error) {
	// Get basic analytics data
	summary, err := ms.repo.GetSummary(ctx, date)
	if err != nil {
		return nil, err
	}

	// Create dashboard analytics
	analytics := &domain.DashboardAnalytics{
		WeeklyAttendance: []int{summary.TotalAttendance, 0, 0, 0, 0, 0, 0}, // Placeholder
		LeaveDistribution: map[string]int{
			"pending":  summary.PendingLeaves,
			"approved": summary.ApprovedLeaves,
			"rejected": summary.RejectedLeaves,
		},
		GeneratedAt: time.Now(),
	}

	return analytics, nil
}

func (ms *MonitoringService) GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	report := &domain.AttendanceReport{
		PeriodStart: startOfMonth,
		PeriodEnd:   endOfMonth,
	}

	return report, nil
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
		"Total Attendance",
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
		excelHandle.SetCellValue(sheet, fmt.Sprintf("E%d", row), summary.TotalAttendance)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("F%d", row), summary.PendingLeaves)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("G%d", row), summary.ApprovedLeaves)
		excelHandle.SetCellValue(sheet, fmt.Sprintf("H%d", row), summary.RejectedLeaves)
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

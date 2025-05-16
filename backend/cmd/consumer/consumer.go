package consumer

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/bootstrap"
	"github.com/aldotp/employee-attendance-system/internal/adapter/consumer"
)

func RunGenerateReportingAttendance(ctx context.Context) {
	b := bootstrap.NewBootstrap(ctx).BuildConsumerGenerateReportingBootstrap()

	con := consumer.NewConsumer(b)
	con.Start(con.GenerateReport)
}

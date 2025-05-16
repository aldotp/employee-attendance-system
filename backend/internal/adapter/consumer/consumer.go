package consumer

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/bootstrap"
	"github.com/aldotp/employee-attendance-system/internal/adapter/handler/worker"
	"go.uber.org/zap"
)

type Operation func()

type consumer struct {
	bootstrap *bootstrap.Bootstrap
	log       *zap.Logger
}

type Consumer interface {
	Start(Operation ...Operation)
	Stop() error

	GenerateReport()
}

func NewConsumer(b *bootstrap.Bootstrap) Consumer {
	return &consumer{
		bootstrap: b,
		log:       b.Log.With(zap.String("from", "consumer")),
	}
}

func (c *consumer) Start(operations ...Operation) {

	for _, op := range operations {
		op()
	}

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan
	c.Stop()

	c.log.Info("Consumer stopped...")

	os.Exit(0)
}

func (c *consumer) Stop() error {
	time.Sleep(1 * time.Second)

	c.log.Info("Consumer stopped...")
	return nil
}

func (c *consumer) GenerateReport() {
	c.log.Info("Consumer registered...", zap.String("job_name", "save_message"))

	worker.NewReportWorker(c.bootstrap).Run()
}

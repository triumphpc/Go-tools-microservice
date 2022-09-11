package mass_transit_server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AleksK1NG/email-microservice/config"
	task "github.com/AleksK1NG/email-microservice/internal/mass_transit/delivery/mass_transit"
	"github.com/AleksK1NG/email-microservice/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

type Server struct {
	db       *sqlx.DB
	amqpConn *amqp.Connection
	logger   logger.Logger
	cfg      *config.Config
}

// NewServer Server constructor
func NewServer(amqpConn *amqp.Connection, logger logger.Logger, cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{amqpConn: amqpConn, logger: logger, cfg: cfg, db: db}
}

func (s *Server) Run() (err error) {
	//metric, err := metrics.CreateMetrics(s.cfg.MassTransitMetrics.URL, s.cfg.MassTransitMetrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.MassTransitMetrics.URL,
		s.cfg.MassTransitMetrics.ServiceName,
	)

	publisher, err := task.NewMassTransitPublisher(s.cfg, s.logger)
	if err != nil {
		return err
	}
	defer publisher.CloseChan()
	s.logger.Info("Mass Transit Publisher initialized")

	//im := interceptors.NewInterceptorManager(s.logger, s.cfg, metric)

	consumer := task.NewConsumer(s.amqpConn, s.logger)

	ctx, cancel := context.WithCancel(context.Background())

	router := echo.New()
	router.GET("/mass_transit_metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(s.cfg.MassTransitMetrics.URL); err != nil {
			s.logger.Errorf("Mass Transit router.Start metrics: %v", err)
			cancel()
		}
	}()

	go func() {
		err := consumer.StartConsumer(s.cfg.RabbitMT)
		if err != nil {
			s.logger.Errorf("StartConsumer for Mass Transit: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	if err := router.Shutdown(ctx); err != nil {
		s.logger.Errorf("Metrics router.Shutdown: %v", err)
	}
	//server.GracefulStop()
	//s.logger.Info("Server Exited Properly")

	return nil

}

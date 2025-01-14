package main

import (
	"log"
	"os"

	"github.com/AleksK1NG/email-microservice/config"
	"github.com/AleksK1NG/email-microservice/internal/mass_transit_server"
	"github.com/AleksK1NG/email-microservice/pkg/jaeger"
	"github.com/AleksK1NG/email-microservice/pkg/logger"
	"github.com/AleksK1NG/email-microservice/pkg/postgres"
	"github.com/AleksK1NG/email-microservice/pkg/rabbitmq"
	"github.com/opentracing/opentracing-go"
)

// Server for mass transit messages
// gRPC service -> RabbitMQ -> Goroutine parallel patterns  -> pg

func main() {
	log.Println("Starting mass transit server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.MassTransitServer.AppVersion,
		cfg.Logger.Level,
		cfg.MassTransitServer.Mode,
		cfg.MassTransitServer.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.MassTransitServer.AppVersion)
	appLogger.Infof("TEST: %#v", cfg.RabbitMT.RepeatExchange)
	appLogger.Infof("TEST: %#v", cfg.RabbitMQ.Exchange)

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		appLogger.Fatal(err)
	}
	defer amqpConn.Close()

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer psqlDB.Close()

	appLogger.Infof("PostgreSQL connected: %#v", psqlDB.Stats())

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	s := mass_transit_server.NewServer(amqpConn, appLogger, cfg, psqlDB)

	appLogger.Fatal(s.Run())
}

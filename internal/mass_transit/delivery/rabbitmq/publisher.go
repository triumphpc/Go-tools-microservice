package task

import (
	"time"

	"github.com/AleksK1NG/email-microservice/config"
	"github.com/AleksK1NG/email-microservice/pkg/logger"
	"github.com/AleksK1NG/email-microservice/pkg/rabbitmq"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type MassTransitPublisher struct {
	amqpChan *amqp.Channel
	cfg      *config.Config
	logger   logger.Logger
}

func NewMassTransitPublisher(cfg *config.Config, logger logger.Logger) (*MassTransitPublisher, error) {
	mqConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		return nil, err
	}
	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Mass transit p.amqpConn.Channel")
	}

	return &MassTransitPublisher{cfg: cfg, logger: logger, amqpChan: amqpChan}, nil
}

// CloseChan Close messages chan
func (p *MassTransitPublisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		p.logger.Errorf("Mass Transit CloseChan: %v", err)
	}
}

// Publish message
func (p *MassTransitPublisher) Publish(body []byte) error {

	p.logger.Infof("Publishing message Work Exchange: %s, RoutingKey: %s", p.cfg.RabbitMT.WorkExchange, p.cfg.RabbitMT.RoutingKey)

	//todo
	if err := p.amqpChan.Publish(
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	publishedMessages.Inc()
	return nil
}

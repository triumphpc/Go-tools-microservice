package task

import (
	"context"

	"github.com/AleksK1NG/email-microservice/config"
	"github.com/AleksK1NG/email-microservice/pkg/logger"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/streadway/amqp"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false
	repeatTTL       = 300

	headerMessageTTL = "x-message-ttl"
	headerDeadLetter = "x-dead-letter-exchange"

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_incoming_rabbitmq_messages_total",
		Help: "The total number of incoming RabbitMQ messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_success_incoming_rabbitmq_messages_total",
		Help: "The total number of success incoming success RabbitMQ messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_error_incoming_rabbitmq_message_total",
		Help: "The total number of error incoming success RabbitMQ messages",
	})
)

// Consumer Rabbitmq consumer
type Consumer struct {
	amqpConn *amqp.Connection
	logger   logger.Logger
}

// NewConsumer Messages Consumer constructor
func NewConsumer(amqpConn *amqp.Connection, logger logger.Logger) *Consumer {
	return &Consumer{amqpConn: amqpConn, logger: logger}
}

// StartConsumer Start new rabbitmq consumer
func (c *Consumer) StartConsumer(cnf config.RabbitMT) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(cnf)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		cnf.WorkQueue,
		"",
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Mass Transit Consume")
	}

	for i := 0; i < cnf.WorkerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Errorf("ch.NotifyClose: %v", chanErr)
	return chanErr
}

// CreateChannel Consume messages
func (c *Consumer) CreateChannel(cnf config.RabbitMT) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error Mass Transit amqpConn.Channel")
	}

	c.logger.Infof("Declaring WorkExchange: %s", cnf.WorkExchange)
	if err = ch.ExchangeDeclare(
		cnf.WorkExchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	); err != nil {
		return nil, errors.Wrap(err, "Error Mass Transit WorkExchange ch.ExchangeDeclare")
	}

	c.logger.Infof("Declaring RepeatExchange: %s", cnf.RepeatExchange)
	if err = ch.ExchangeDeclare(
		cnf.RepeatExchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	); err != nil {
		return nil, errors.Wrap(err, "Error Mass Transit RepeatExchange ch.ExchangeDeclare")
	}

	c.logger.Infof("Declaring DeadExchange: %s", cnf.RepeatExchange)
	if err = ch.ExchangeDeclare(
		cnf.DeadExchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	); err != nil {
		return nil, errors.Wrap(err, "Error Mass Transit DeadExchange ch.ExchangeDeclare")
	}

	c.logger.Infof("Declaring WorkQueue: %s", cnf.WorkQueue)
	workQueue, err := ch.QueueDeclare(
		cnf.WorkQueue,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		amqp.Table{
			headerDeadLetter: cnf.RepeatExchange,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error WorkQueue ch.QueueDeclare")
	}

	c.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, RoutingKey: %v",
		workQueue.Name,
		workQueue.Messages,
		workQueue.Consumers,
		cnf.WorkExchange,
		cnf.RoutingKey,
	)

	err = ch.QueueBind(
		workQueue.Name,
		cnf.RoutingKey,
		cnf.WorkExchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error WorkQueue ch.QueueBind")
	}

	c.logger.Infof("Declaring RepeatQueue: %s", cnf.WorkQueue)
	repeatQueue, err := ch.QueueDeclare(
		cnf.RepeatQueue,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		amqp.Table{
			headerDeadLetter: cnf.WorkExchange,
			headerMessageTTL: repeatTTL,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error repeatQueue ch.QueueDeclare")
	}

	c.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v",
		repeatQueue.Name,
		repeatQueue.Messages,
		repeatQueue.Consumers,
		cnf.RepeatExchange,
	)

	err = ch.QueueBind(
		repeatQueue.Name,
		cnf.RoutingKey,
		cnf.RepeatExchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error RepeatQueue ch.QueueBind")
	}

	c.logger.Infof("Declaring DeadQueue: %s", cnf.WorkQueue)
	deadQueue, err := ch.QueueDeclare(
		cnf.DeadQueue,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		amqp.Table{
			headerDeadLetter: cnf.WorkExchange,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error deadQueue ch.QueueDeclare")
	}

	c.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v",
		deadQueue.Name,
		deadQueue.Messages,
		deadQueue.Consumers,
		cnf.DeadExchange,
	)

	err = ch.QueueBind(
		deadQueue.Name,
		cnf.RoutingKey,
		cnf.DeadExchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error deadQueue ch.QueueBind")
	}

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error Mass Transit  ch.Qos")
	}

	return ch, nil
}

func (c *Consumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {

	for delivery := range messages {
		//span, ctx := opentracing.StartSpanFromContext(ctx, "Mass Transit Consumer.worker")

		c.logger.Infof("processDeliveries deliveryTag% v", delivery.DeliveryTag)

		incomingMessages.Inc()

		//err := c.emailUC.SendEmail(ctx, delivery.Body)
		//if err != nil {
		//	if err := delivery.Reject(false); err != nil {
		//		c.logger.Errorf("Err delivery.Reject: %v", err)
		//	}
		//	c.logger.Errorf("Failed to process delivery: %v", err)
		//	errorMessages.Inc()
		//} else {
		err := delivery.Ack(false)
		if err != nil {
			c.logger.Errorf("Failed to acknowledge delivery: %v", err)
		}
		//	successMessages.Inc()
		//}
		//span.Finish()
	}

	c.logger.Info("Deliveries channel closed")
}

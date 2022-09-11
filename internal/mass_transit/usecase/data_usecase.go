package usecase

import (
	"context"
	"fmt"

	"github.com/AleksK1NG/email-microservice/config"
	"github.com/AleksK1NG/email-microservice/internal/email"
	task "github.com/AleksK1NG/email-microservice/internal/mass_transit/delivery/mass_transit"
	"github.com/AleksK1NG/email-microservice/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

// DataUseCase Data cases
type DataUseCase struct {
	repo      email.EmailsRepository
	logger    logger.Logger
	cfg       *config.Config
	publisher task.MassTransitPublisher
}

// NewDataUseCase Image useCase constructor
func NewDataUseCase(repo email.EmailsRepository, logger logger.Logger, cfg *config.Config, publisher task.MassTransitPublisher) *DataUseCase {
	return &DataUseCase{repo: repo, logger: logger, cfg: cfg, publisher: publisher}
}

// PublishDataToQueue Publish data to rabbitmq
func (e *DataUseCase) PublishDataToQueue(ctx context.Context, data interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DataUseCase.PublishEmailToQueue")
	defer span.Finish()

	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("can't convert data to bytes")
	}

	return e.publisher.Publish(buf)
}

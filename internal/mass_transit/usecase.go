package mass_transit

import (
	"context"
)

// UseCase  interface
type UseCase interface {
	// PublishDataToQueue - publish data to queue
	PublishDataToQueue(ctx context.Context, data interface{}) error
}

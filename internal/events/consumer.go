package events

import "context"

type Consumer interface {
	Consume(ctx context.Context) error
	CloseConnection() error
}

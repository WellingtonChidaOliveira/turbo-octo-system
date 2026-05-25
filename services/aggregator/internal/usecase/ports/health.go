package ports

import "context"

type QueueHealthChecker interface {
	CheckQueue(ctx context.Context) error
}

type StorageHealthChecker interface {
	CheckStorage(ctx context.Context) error
}

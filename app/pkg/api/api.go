package api

import (
	"context"
)

type (
	// AvailabilityChecker - проверяет доступность объекта по его ID.
	AvailabilityChecker interface {
		CheckAvailability(ctx context.Context, itemID uint64) error
	}
)

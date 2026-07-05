package module

import (
	"github.com/mondegor/go-core/errors"
)

var (
	// ErrContainerNotFound - container with ID not found.
	ErrContainerNotFound = errors.NewUserProto("ContainerNotFound", "container with ID={Id} not found")

	// ErrStoreNotFound - store with ID not found.
	ErrStoreNotFound = errors.NewUserProto("StoreNotFound", "store with ID={Id} not found")

	// ErrStockNotFound - stock with ID not found.
	ErrStockNotFound = errors.NewUserProto("StockNotFound", "stock with ID={Id} not found")

	// ErrLocationIsOccupied - location is occupied.
	ErrLocationIsOccupied = errors.New("location is occupied")
)

package usecase

import (
	"context"

	"github.com/mondegor/go-core/errors"

	"print-shop-back/internal/warehousing/enum/locationkind"
)

const (
	locationCount = locationkind.Group + 1
)

type (
	// RefreshLocations - comment struct.
	RefreshLocations struct {
		useCaseRefreshLocations [locationCount]refreshLocationUseCase
		errorWrapper            errors.Wrapper
	}

	refreshLocationUseCase interface {
		Execute(ctx context.Context, storeIDs []uint64) error
	}
)

// NewRefreshLocations - создаёт объект RefreshLocations.
func NewRefreshLocations(
	useCaseRefreshStores refreshLocationUseCase,
	useCaseRefreshGroupContainers refreshLocationUseCase,
) *RefreshLocations {
	rl := &RefreshLocations{
		errorWrapper: errors.NewServiceOperationFailedWrapper(),
	}

	rl.useCaseRefreshLocations[locationkind.Store] = useCaseRefreshStores
	rl.useCaseRefreshLocations[locationkind.Group] = useCaseRefreshGroupContainers

	return rl
}

// Execute - comment method.
func (uc *RefreshLocations) Execute(ctx context.Context, locationIDs []uint64) error {
	var (
		ns [locationCount]int
		nl [locationCount][]uint64
	)

	for i := range locationIDs {
		kind := locationkind.ByID(locationIDs[i])
		if kind < locationCount {
			ns[kind]++
		}
	}

	for kind := range ns {
		if ns[kind] > 0 {
			nl[kind] = make([]uint64, 0, ns[kind])
		}
	}

	for i := range locationIDs {
		kind := locationkind.ByID(locationIDs[i])
		if kind < locationCount {
			nl[kind] = append(nl[kind], locationIDs[i])
		}
	}

	// TODO: можно сделать в отдельных горутинах
	for i := range nl {
		if len(nl[i]) == 0 {
			continue
		}

		if err := uc.useCaseRefreshLocations[i].Execute(ctx, nl[i]); err != nil {
			return uc.errorWrapper.Wrap(err)
		}
	}

	return nil
}

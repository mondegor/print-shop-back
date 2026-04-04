package usecase

import (
	"context"
	"math"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/slices/suint64"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/util/locationstock"
	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
	"github.com/mondegor/print-shop-back/internal/warehousing/module"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// RefreshStoreContainersVolume - comment struct.
	RefreshStoreContainersVolume struct {
		storeStorage storeStorage
		storageStock stockStorage
		logger       mrlog.Logger
		errorWrapper errors.Wrapper
	}

	storeStorage interface {
		UpdateContainersVolume(ctx context.Context, rows []dto.LocationContainersVolume) error
	}
)

// NewRefreshStoreContainersVolume - создаёт объект RefreshStoreContainersVolume.
func NewRefreshStoreContainersVolume(
	storageStore storeStorage,
	storageStock stockStorage,
	logger mrlog.Logger,
) *RefreshStoreContainersVolume {
	return &RefreshStoreContainersVolume{
		storeStorage: storageStore,
		storageStock: storageStock,
		logger:       logger,
		errorWrapper: errors.NewServiceOperationFailedWrapper(),
	}
}

// Execute - comment method.
func (uc *RefreshStoreContainersVolume) Execute(ctx context.Context, storeIDs []uint64) error {
	storeIDs = suint64.FilterFunc(
		storeIDs,
		func(el uint64) bool {
			return locationkind.Is(el, locationkind.Store)
		},
	)

	// TODO: склады разбить на чанки, чтобы уменьшить кол-во получаемых стоков
	if len(storeIDs) == 0 {
		return nil
	}

	storeIDs = suint64.SortedUnique(storeIDs)

	stocks, hasNext, err := uc.storageStock.FetchByLocationIDs(
		ctx,
		storeIDs,
		xtype.StockCursor{
			Limit: min(module.GroupContainersMax*len(storeIDs), math.MaxInt16),
		},
	)
	if err != nil {
		return uc.errorWrapper.Wrap(err)
	}

	if hasNext {
		uc.logger.Error(
			ctx,
			"there are more containers in stores than allowed value",
			"store_ids", storeIDs,
			"limit", len(stocks),
		)
	}

	stores := make([]dto.LocationContainersVolume, 0, len(storeIDs))

	for s := range locationstock.ChunkByLocationID(stocks) {
		stores = append(
			stores,
			dto.LocationContainersVolume{
				LocationID:  s[0].LocationID, // существование гарантирует метод ChunkByLocationID()
				TotalVolume: locationstock.CalcContainersVolume(s),
			},
		)
	}

	// помечаются все склады со стоками обработанными
	for i := range stocks {
		if index := suint64.BinaryIndex(storeIDs, stocks[i].LocationID); index >= 0 {
			storeIDs[index] = 0
		}
	}

	for _, storeID := range storeIDs {
		if storeID == 0 {
			continue
		}

		// все необработанные склады не содержат стоков
		stores = append(
			stores,
			dto.LocationContainersVolume{
				LocationID: storeID,
			},
		)
	}

	return uc.storeStorage.UpdateContainersVolume(ctx, stores)
}

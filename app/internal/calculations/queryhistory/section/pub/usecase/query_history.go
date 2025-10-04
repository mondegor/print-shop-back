package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/entity"
)

type (
	// QueryHistory - comment struct.
	QueryHistory struct {
		storage      pub.QueryResultStorage
		eventEmitter mrevent.Emitter
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewQueryHistory - создаёт объект QueryHistory.
func NewQueryHistory(
	storage pub.QueryResultStorage,
	eventEmitter mrevent.Emitter,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *QueryHistory {
	return &QueryHistory{
		storage:      storage,
		eventEmitter: mrevent.NewSourceEmitter(eventEmitter, entity.ModelNameQueryHistory),
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameQueryHistory),
	}
}

// GetItem - comment method.
func (uc *QueryHistory) GetItem(ctx context.Context, itemID uuid.UUID) (entity.QueryHistoryItem, error) {
	if itemID == uuid.Nil {
		return entity.QueryHistoryItem{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.QueryHistoryItem{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	// обновление счётчика посещений
	// TODO: send to queue
	go func() {
		// if err := uc.storage.UpdateQuantity(ctx, itemID); err != nil {
		//	 mrlog.Ctx(ctx).Error().Err(err).Send()
		// }
	}()

	return item, nil
}

// Create - comment method.
func (uc *QueryHistory) Create(ctx context.Context, item entity.QueryHistoryItem) (itemID uuid.UUID, err error) {
	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrargs.Group{"id": itemID})

	return itemID, nil
}

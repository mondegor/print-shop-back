package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrevent"

	"print-shop-back/internal/calculations/queryhistory/section/pub"
	"print-shop-back/internal/calculations/queryhistory/section/pub/entity"
)

type (
	// QueryHistory - comment struct.
	QueryHistory struct {
		storage      pub.QueryResultStorage
		eventEmitter mrevent.Emitter
		errorWrapper errors.Wrapper
	}
)

// NewQueryHistory - создаёт объект QueryHistory.
func NewQueryHistory(
	storage pub.QueryResultStorage,
	eventEmitter mrevent.Emitter,
) *QueryHistory {
	return &QueryHistory{
		storage:      storage,
		eventEmitter: mrevent.EmitterWithSource(eventEmitter, entity.ModelNameQueryHistory),
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetItem - comment method.
func (uc *QueryHistory) GetItem(ctx context.Context, itemID uuid.UUID) (entity.QueryHistoryItem, error) {
	if itemID == uuid.Nil {
		return entity.QueryHistoryItem{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.QueryHistoryItem{}, uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	// обновление счётчика посещений
	// TODO: send to queue
	go func() {
		// if err := uc.storage.UpdateQuantity(ctx, itemID); err != nil {
		//	 log.Ctx(ctx).Error().Err(err).Send()
		// }
	}()

	return item, nil
}

// Create - comment method.
func (uc *QueryHistory) Create(ctx context.Context, item entity.QueryHistoryItem) (itemID uuid.UUID, err error) {
	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", "itemId", itemID)

	return itemID, nil
}

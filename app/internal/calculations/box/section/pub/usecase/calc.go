package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/entity"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CalcResult - comment struct.
	CalcResult struct {
		storage      CalcResultStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewBox - создаёт объект CalcResult.
func NewBox(storage CalcResultStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *CalcResult {
	return &CalcResult{
		storage:      storage,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetItem - comment method.
func (uc *CalcResult) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.CalcResult, error) {
	if itemID < 1 {
		return entity.CalcResult{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.CalcResult{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCalcResult, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *CalcResult) Create(ctx context.Context, item entity.CalcResult) (mrtype.KeyInt32, error) {
	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCalcResult)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *CalcResult) Store(ctx context.Context, item entity.CalcResult) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	tagVersion, err := uc.storage.Insert(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCalcResult)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *CalcResult) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCalcResult,
		data,
	)
}

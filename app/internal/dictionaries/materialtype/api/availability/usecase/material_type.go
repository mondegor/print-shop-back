package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      availability.MaterialTypeStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(storage availability.MaterialTypeStorage, errorWrapper mrcore.UseCaseErrorWrapper) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *MaterialType) CheckingAvailability(ctx context.Context, itemID uint64) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID == 0 {
		return api.ErrMaterialTypeRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrMaterialTypeNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.MaterialTypeAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrMaterialTypeNotAvailable.New(itemID)
	}

	return nil
}

func (uc *MaterialType) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.MaterialTypeAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}

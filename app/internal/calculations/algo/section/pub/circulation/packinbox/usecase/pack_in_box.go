package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

type (
	// CirculationPackInBox - comment struct.
	CirculationPackInBox struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewCirculationPackInBox - создаёт объект CirculationPackInBox.
func NewCirculationPackInBox(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *CirculationPackInBox {
	return &CirculationPackInBox{
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *CirculationPackInBox) CalcQuantity(ctx context.Context, raw entity.RawData) (entity.AlgoResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.AlgoResult{}, err
	}

	uc.emitEvent(ctx, "CalcQuantity", mrmsg.Data{"raw": parsedData})

	return entity.AlgoResult{
		Format: parsedData.Format,
	}, nil
}

func (uc *CirculationPackInBox) parse(data entity.RawData) (entity.ParsedData, error) {
	format, err := rect.ParseFormat(data.Format)
	if err != nil {
		return entity.ParsedData{}, err
	}

	return entity.ParsedData{
		Format: format,
	}, nil
}

func (uc *CirculationPackInBox) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNamePackInBox,
		data,
	)
}

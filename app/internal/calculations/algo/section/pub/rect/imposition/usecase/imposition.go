package usecase

import (
	"context"
	"errors"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// RectImposition - comment struct.
	RectImposition struct {
		algo         *imposition.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewRectImposition - создаёт объект RectImposition.
func NewRectImposition(algo *imposition.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *RectImposition {
	return &RectImposition{
		algo:         algo,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *RectImposition) Calc(ctx context.Context, raw entity.RawData) (imposition.AlgoResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return imposition.AlgoResult{}, err
	}

	result, err := uc.algo.Calc(parsedData.Item, parsedData.Out, parsedData.Opts)
	if err != nil {
		return imposition.AlgoResult{}, err
	}

	if len(result.Fragments) == 0 {
		return imposition.AlgoResult{}, errors.New("error") // TODO: result NULL
	}

	uc.emitEvent(ctx, "Calc", mrmsg.Data{"raw": parsedData})

	return result, nil
}

func (uc *RectImposition) parse(data entity.RawData) (entity.ParsedData, error) {
	itemFormat, err := rect.ParseFormat(data.ItemFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemFormat error
	}

	itemBorderFormat := rect.Format{} // optional

	if data.ItemBorderFormat != "" {
		itemBorderFormat, err = rect.ParseFormat(data.ItemBorderFormat)
		if err != nil {
			return entity.ParsedData{}, err // TODO: itemBorderFormat error
		}
	}

	outFormat, err := rect.ParseFormat(data.OutFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: outFormat error
	}

	return entity.ParsedData{
		Item: rect.Item{
			Format: itemFormat,
			Border: itemBorderFormat,
		},
		Out: outFormat,
		Opts: imposition.Options{
			AllowRotation: data.AllowRotation,
			UseMirror:     data.UseMirror,
		},
	}, nil
}

func (uc *RectImposition) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameRectImposition,
		data,
	)
}

package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

func (ht *RectInsideOutside) parseRequest(data CalcRectInsideOutsideRequest) (entity.ParsedData, error) {
	inFormat, err := rect.ParseFormat(data.InFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemInFormat error
	}

	outFormat, err := rect.ParseFormat(data.OutFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemOutFormat error
	}

	return entity.ParsedData{
		In:  inFormat.Transform(measure.OneThousandth),  // mm -> m
		Out: outFormat.Transform(measure.OneThousandth), // mm -> m
	}, nil
}

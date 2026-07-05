package httpv1

import (
	"print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1/model"
	"print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/dto"
	"print-shop-back/pkg/mrcalc/measure"
	"print-shop-back/pkg/mrcalc/s2/rect2d"
)

func (ht *SheetInsideOutside) parseRequest(data model.SheetInsideOutsideQuantityRequest) (dto.ParsedData, error) {
	inFormat, err := rect2d.ParseFormat(data.InFormat)
	if err != nil {
		return dto.ParsedData{}, err // TODO: itemInFormat error
	}

	outFormat, err := rect2d.ParseFormat(data.OutFormat)
	if err != nil {
		return dto.ParsedData{}, err // TODO: itemOutFormat error
	}

	return dto.ParsedData{
		In:  inFormat.Transform(measure.OneThousandth),  // mm -> m
		Out: outFormat.Transform(measure.OneThousandth), // mm -> m
	}, nil
}

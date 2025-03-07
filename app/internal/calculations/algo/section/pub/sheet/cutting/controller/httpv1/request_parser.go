package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func (ht *SheetCutting) parseRequest(data model.SheetCuttingQuantityRequest) (dto.ParsedData, error) {
	fragments := make([]rect2d.Layout, 0, len(data.Layouts))

	for _, str := range data.Layouts {
		fragment, err := rect2d.ParseLayout(str)
		if err != nil {
			return dto.ParsedData{}, err // TODO: itemInFormat error
		}

		fragments = append(fragments, fragment)
	}

	distanceFormat, err := rect2d.ParseFormat(data.DistanceFormat)
	if err != nil {
		return dto.ParsedData{}, err // TODO: distanceFormat error
	}

	return dto.ParsedData{
		Fragments:      fragments,
		DistanceFormat: distanceFormat.Transform(measure.OneThousandth), // mm -> m
	}, nil
}

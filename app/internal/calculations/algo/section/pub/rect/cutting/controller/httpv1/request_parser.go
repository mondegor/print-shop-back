package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

func (ht *RectCutting) parseRequest(data CalcCuttingQuantityRequest) (entity.ParsedData, error) {
	fragments := make([]base.Fragment, 0, len(data.Fragments))

	for _, str := range data.Fragments {
		fragment, err := base.ParseFragment(str)
		if err != nil {
			return entity.ParsedData{}, err // TODO: itemInFormat error
		}

		fragments = append(fragments, fragment)
	}

	distanceFormat, err := rect.ParseFormat(data.DistanceFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: distanceFormat error
	}

	return entity.ParsedData{
		Fragments:      fragments,
		DistanceFormat: distanceFormat.Transform(measure.OneThousandth), // mm -> m
	}, nil
}

package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

func (ht *RectImposition) parseRequest(data CalcRectImpositionRequest) (entity.ParsedData, error) {
	itemFormat, err := rect.ParseFormat(data.ItemFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemFormat error
	}

	itemDistance := rect.Format{} // optional

	if data.ItemDistance != "" {
		itemDistance, err = rect.ParseFormat(data.ItemDistance)
		if err != nil {
			return entity.ParsedData{}, err // TODO: itemDistance error
		}
	}

	outFormat, err := rect.ParseFormat(data.OutFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: outFormat error
	}

	return entity.ParsedData{
		Item: rect.Item{
			Format:   itemFormat.Transform(measure.OneThousandth),   // mm -> m
			Distance: itemDistance.Transform(measure.OneThousandth), // mm -> m
		},
		Out: outFormat.Transform(measure.OneThousandth), // mm -> m
		Opts: imposition.Options{
			AllowRotation: !data.DisableRotation,
			UseMirror:     data.UseMirror,
		},
	}, nil
}

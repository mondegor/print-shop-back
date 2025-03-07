package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func (ht *SheetImposition) parseRequest(data model.SheetImpositionRequest) (dto.ParsedData, error) {
	itemFormat, err := rect2d.ParseFormat(data.ItemFormat)
	if err != nil {
		return dto.ParsedData{}, err // TODO: itemFormat error
	}

	itemDistance := rect2d.Format{} // optional

	if data.ItemDistance != "" {
		itemDistance, err = rect2d.ParseFormat(data.ItemDistance)
		if err != nil {
			return dto.ParsedData{}, err // TODO: itemDistance error
		}
	}

	outFormat, err := rect2d.ParseFormat(data.OutFormat)
	if err != nil {
		return dto.ParsedData{}, err // TODO: outFormat error
	}

	return dto.ParsedData{
		Element:  itemFormat.Transform(measure.OneThousandth),   // mm -> m
		Distance: itemDistance.Transform(measure.OneThousandth), // mm -> m
		Out:      outFormat.Transform(measure.OneThousandth),    // mm -> m
		Opts: imposition.Options{
			AllowRotation: !data.DisableRotation,
			UseMirror:     data.UseMirror,
		},
	}, nil
}

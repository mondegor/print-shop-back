package httpv1

import (
	httpmodel "github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func (ht *PackInStack) parseRequest(data httpmodel.SheetPackInStackRequest) (dto.ParsedData, error) {
	sheetFormat, err := rect2d.ParseFormat(data.Sheet.Format)
	if err != nil {
		return dto.ParsedData{}, err // TODO: sheetFormat error
	}

	return dto.ParsedData{
		SheetHeap: model.SheetStack{
			Sheet: model.Sheet{
				Format:    sheetFormat.Transform(measure.OneThousandth),         // mm -> m
				Thickness: float64(data.Sheet.Thickness) * measure.OneMillionth, // mkm -> m
				Density:   float64(data.Sheet.Density) * measure.OneThousandth,  // g/m2 -> kg/m2
			},
			Quantity: data.Sheet.Quantity,
		},
		QuantityInStack: data.QuantityInStack,
	}, nil
}

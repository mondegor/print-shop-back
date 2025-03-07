package httpv1

import (
	httpmodel "github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3/rect3d"
)

func (ht *BoxPackInBox) parseRequest(data httpmodel.CalcBoxPackInBoxRequest) (dto.ParsedData, error) {
	productFormat, err := rect3d.ParseFormat(data.Product.Format)
	if err != nil {
		return dto.ParsedData{}, err // TODO: productFormat error
	}

	boxFormat, err := rect3d.ParseFormat(data.Box.Format)
	if err != nil {
		return dto.ParsedData{}, err // TODO: boxFormat error
	}

	boxMargins, err := rect3d.ParseFormat(data.Box.Margins)
	if err != nil {
		return dto.ParsedData{}, err // TODO: boxMargins error
	}

	// TODO: if boxMargins > boxFormat

	return dto.ParsedData{
		ProductHeap: model.ProductStack{
			Product: model.Product{
				Format: productFormat.Transform(measure.OneThousandth),       // mm -> m
				Weight: float64(data.Product.Weight) * measure.OneThousandth, // g -> kg
			},
			Quantity: data.Product.Quantity,
		},
		Box: model.Box{
			Format:    boxFormat.Transform(measure.OneThousandth),          // mm -> m
			Thickness: float64(data.Box.Thickness) * measure.OneMillionth,  // mkm -> m
			Margins:   boxMargins.Transform(measure.OneThousandth),         // mm -> m
			Weight:    float64(data.Box.Weight) * measure.OneThousandth,    // g -> kg
			MaxWeight: float64(data.Box.MaxWeight) * measure.OneThousandth, // g -> kg
		},
	}, nil
}

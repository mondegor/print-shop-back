package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/packinbox"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/parallelepiped"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

func (ht *PackInBox) parseRequest(data CalcCirculationPackInBoxRequest) (entity.ParsedData, error) {
	productFormat, err := rect.ParseFormat(data.Product.Format)
	if err != nil {
		return entity.ParsedData{}, err // TODO: productFormat error
	}

	boxFormat, err := parallelepiped.ParseFormat(data.Box.Format)
	if err != nil {
		return entity.ParsedData{}, err // TODO: boxFormat error
	}

	boxMargins, err := parallelepiped.ParseFormat(data.Box.Margins)
	if err != nil {
		return entity.ParsedData{}, err // TODO: boxMargins error
	}

	// TODO: if boxMargins > boxFormat

	return entity.ParsedData{
		Product: packinbox.Product{
			Format:    productFormat.Transform(measure.OneThousandth),         // mm -> m
			Thickness: float64(data.Product.Thickness) * measure.OneMillionth, // mkm -> m
			WeightM2:  float64(data.Product.WeightM2) * measure.OneThousandth, // g/m2 -> kg/m2
			Quantity:  data.Product.Quantity,
		},
		Box: packinbox.Box{
			Format:    boxFormat.Transform(measure.OneThousandth),          // mm -> m
			Thickness: float64(data.Box.Thickness) * measure.OneMillionth,  // mkm -> m
			Margins:   boxMargins.Transform(measure.OneThousandth),         // mm -> m
			Weight:    float64(data.Box.Weight) * measure.OneThousandth,    // g -> kg
			MaxWeight: float64(data.Box.MaxWeight) * measure.OneThousandth, // g -> kg
		},
	}, nil
}

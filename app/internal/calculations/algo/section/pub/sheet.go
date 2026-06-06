package pub

import (
	"context"

	cuttingmodel "print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1/model"
	cuttingdto "print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/dto"
	impositionmodel "print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1/model"
	impositiondto "print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/dto"
	insideoutsidemodel "print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1/model"
	insideoutsidedto "print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/dto"
	packinstackmodel "print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1/model"
	packinstackdto "print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/dto"
)

type (
	// SheetCuttingUseCase - comment interface.
	SheetCuttingUseCase interface {
		CalcQuantity(ctx context.Context, data cuttingdto.ParsedData) (cuttingmodel.SheetCuttingQuantityResult, error)
	}

	// SheetImpositionUseCase - comment interface.
	SheetImpositionUseCase interface {
		Calc(ctx context.Context, data impositiondto.ParsedData) (impositionmodel.SheetImpositionResponse, error)
		CalcVariants(ctx context.Context, data impositiondto.ParsedData) (impositionmodel.SheetImpositionVariantsResponse, error)
	}

	// SheetInsideOutsideUseCase - comment interface.
	SheetInsideOutsideUseCase interface {
		CalcQuantity(ctx context.Context, data insideoutsidedto.ParsedData) (insideoutsidemodel.SheetInsideOutsideQuantityResponse, error)
		CalcMax(ctx context.Context, data insideoutsidedto.ParsedData) (insideoutsidemodel.SheetInsideOutsideMaxResponse, error)
	}

	// SheetPackInStackUseCase - comment interface.
	SheetPackInStackUseCase interface {
		Calc(ctx context.Context, data packinstackdto.ParsedData) (packinstackmodel.SheetPackInStackResponse, error)
	}
)

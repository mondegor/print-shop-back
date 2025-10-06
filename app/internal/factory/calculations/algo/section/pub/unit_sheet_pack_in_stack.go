package pub

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/usecase"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/packinstack"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initSheetPackInStackController(
	eventEmitter mrevent.Emitter,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	packInStackAlgoSheet := packinstack.New()

	useCase := usecase.NewSheetPackInStack(packInStackAlgoSheet, eventEmitter)

	controller := httpv1.NewSheetPackInStack(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

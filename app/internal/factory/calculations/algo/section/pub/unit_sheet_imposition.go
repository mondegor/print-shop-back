package pub

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/usecase"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initSheetImpositionController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	algoComponent := imposition.New(logger)

	useCase := usecase.NewSheetImposition(algoComponent, eventEmitter)

	controller := httpv1.NewSheetImposition(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

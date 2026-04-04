package pub

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/usecase"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initBoxSheetCuttingController(
	eventEmitter mrevent.Emitter,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	useCase := usecase.NewSheetCutting(eventEmitter)

	controller := httpv1.NewSheetCutting(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

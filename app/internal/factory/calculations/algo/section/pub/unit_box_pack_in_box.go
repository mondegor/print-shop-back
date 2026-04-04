package pub

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/usecase"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/box/packinbox"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initBoxPackInBoxController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	impAlgo := imposition.New(logger)
	packInBoxAlgo := packinbox.New(impAlgo)

	useCase := usecase.NewBoxPackInBox(packInBoxAlgo, logger, eventEmitter)

	controller := httpv1.NewBoxPackInBox(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

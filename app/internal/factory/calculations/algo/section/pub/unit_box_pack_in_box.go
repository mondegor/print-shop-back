package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/box/packinbox"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
)

func createUnitBoxPackInBox(opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBoxPackInBox(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBoxPackInBox(opts algo.Options) (*httpv1.BoxPackInBox, error) { //nolint:unparam
	impAlgo := imposition.New(opts.Logger)
	packInBoxAlgo := packinbox.New(impAlgo)

	useCase := usecase.NewBoxPackInBox(packInBoxAlgo, opts.Logger, opts.EventEmitter)
	controller := httpv1.NewBoxPackInBox(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

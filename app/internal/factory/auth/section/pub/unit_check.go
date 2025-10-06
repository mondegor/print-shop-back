package pub

import (
	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

func createUnitCheck(opts auth.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCheck(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

//nolint:unparam
func newUnitCheck(opts auth.Options) (*httpv1.Check, error) {
	useCase := check.NewAuthHelper(
		createCheckUserPostgres(opts),
		createUserRealmPostgres(opts),
		contactaddress.NewParser(), // ??????
		opts.UseCaseErrorWrapper,
	)

	controller := httpv1.NewCheck(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

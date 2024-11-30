package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/usecase"
)

func createUnitCompanyPage(ctx context.Context, opts provideraccounts.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(_ context.Context, opts provideraccounts.Options) (*httpv1.CompanyPage, error) { //nolint:unparam
	storage := repository.NewCompanyPagePostgres(opts.DBConnManager)
	useCase := usecase.NewCompanyPage(storage, opts.UseCaseErrorWrapper)
	controller := httpv1.NewCompanyPage(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
		opts.UnitCompanyPage.LogoURLBuilder,
	)

	return controller, nil
}

package prov

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/usecase"
)

func createUnitCompanyPage(ctx context.Context, opts provideraccounts.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	if c, err := newUnitCompanyPageLogo(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(_ context.Context, opts provideraccounts.Options) (*httpv1.CompanyPage, error) { //nolint:unparam
	storage := repository.NewCompanyPagePostgres(opts.DBConnManager)
	useCase := usecase.NewCompanyPage(
		opts.DBConnManager,
		storage,
		opts.EventEmitter,
		opts.UseCaseHelper,
		opts.UnitCompanyPage.LogoURLBuilder,
	)
	controller := httpv1.NewCompanyPage(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

func newUnitCompanyPageLogo(_ context.Context, opts provideraccounts.Options) (*httpv1.CompanyPageLogo, error) { //nolint:unparam
	storage := repository.NewCompanyPageLogoPostgres(opts.DBConnManager)
	useCase := usecase.NewCompanyPageLogo(
		storage,
		opts.UnitCompanyPage.LogoFileAPI,
		opts.Locker,
		opts.EventEmitter,
		opts.UseCaseHelper,
	)
	controller := httpv1.NewCompanyPageLogo(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

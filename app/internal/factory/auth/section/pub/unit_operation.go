package pub

import (
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/usecase/operation"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

func createUnitOperation(opts auth.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitOperation(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

//nolint:unparam
func newUnitOperation(opts auth.Options) (*httpv1.Operation, error) {
	useCaseConfirmOperation := operation.NewConfirmOperation(
		opts.DBConnManager,
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		secureoperation.NewConfirmCode(
			crypt.NewTokenGenerator(int(opts.OperationConfirm.TokenLength)), // DEFAULT
			crypt.NewCodeGenerator(int(opts.OperationConfirm.CodeLength)),   // DEFAULT
		),
		opts.UseCaseErrorWrapper,
	)

	useCaseResendConfirmCode := operation.NewResendCode(
		opts.DBConnManager,
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		secureoperation.NewResendCode(
			crypt.NewTokenGenerator(int(opts.OperationConfirm.TokenLength)), // DEFAULT
			crypt.NewCodeGenerator(int(opts.OperationConfirm.CodeLength)),   // DEFAULT
		),
		opts.UseCaseErrorWrapper,
	)

	controller := httpv1.NewOperation(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCaseConfirmOperation,
		useCaseResendConfirmCode,
		bag.NewOperationResponse(opts.WithDebugInfo),
	)

	return controller, nil
}

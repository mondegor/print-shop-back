package pub

import (
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-components/mrauth/usecase/operation"
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

func initConfirmOperationUseCase(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	storageSecureOperation *repository.SecureOperationPostgres,
	notifierAPI mrnotifier.NoticeProducer,
	operationConfirm auth.OperationConfirm,
) *operation.ConfirmOperation {
	return operation.NewConfirmOperation(
		dbConnManager,
		storageSecureOperation,
		notifierAPI,
		secureoperation.NewConfirmCode(
			crypt.NewTokenGenerator(int(operationConfirm.TokenLength)), // DEFAULT
			crypt.NewCodeGenerator(int(operationConfirm.CodeLength)),   // DEFAULT
		),
		useCaseErrorWrapper,
	)
}

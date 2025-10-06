package pub

import (
	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initCheckController(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	storageCheckUser *repository.CheckUserPostgres,
	storageUserRealm *repository.UserRealmPostgres,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	useCase := check.NewAuthHelper(
		storageCheckUser,
		storageUserRealm,
		contactaddress.NewParser(), // ??????
		useCaseErrorWrapper,
	)

	controller := httpv1.NewCheck(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

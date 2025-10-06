package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initPaperFactureController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax uint64,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.PaperFacture{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperFacturePostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewPaperFacture(storage, eventEmitter, useCaseErrorWrapper)

	controller := httpv1.NewPaperFacture(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}

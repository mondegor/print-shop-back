package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initPaperColorController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax uint64,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.PaperColor{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperColorPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewPaperColor(storage, eventEmitter, useCaseErrorWrapper)

	controller := httpv1.NewPaperColor(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}

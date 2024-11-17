package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// LaminatePostgres - comment struct.
	LaminatePostgres struct {
		client          mrstorage.DBConnManager
		repoTypeIDs     db.ColumnFetcher[mrenum.ItemStatus, uint64]
		repoThicknesses db.ColumnFetcher[mrenum.ItemStatus, measure.Meter]
	}
)

// NewLaminatePostgres - создаёт объект LaminatePostgres.
func NewLaminatePostgres(client mrstorage.DBConnManager) *LaminatePostgres {
	return &LaminatePostgres{
		client: client,
		repoTypeIDs: db.NewColumnFetcher[mrenum.ItemStatus, uint64](
			client,
			module.DBTableNameLaminates,
			"type_id",
			"laminate_status",
			module.DBFieldDeletedAt,
		),
		repoThicknesses: db.NewColumnFetcher[mrenum.ItemStatus, measure.Meter](
			client,
			module.DBTableNameLaminates,
			"laminate_thickness",
			"laminate_status",
			module.DBFieldDeletedAt,
		),
	}
}

// Fetch - comment method.
func (re *LaminatePostgres) Fetch(ctx context.Context, _ entity.LaminateParams) ([]entity.Laminate, error) {
	sql := `
        SELECT
            laminate_id,
			laminate_article,
			laminate_caption,
			type_id,
			laminate_length,
			laminate_width,
			laminate_thickness,
			laminate_weight_m2
        FROM
            ` + module.DBTableNameLaminates + `
        WHERE
            laminate_status = $1 AND deleted_at IS NULL
        ORDER BY
            laminate_caption ASC, laminate_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.Laminate, 0)

	for cursor.Next() {
		var row entity.Laminate

		err = cursor.Scan(
			&row.ID,
			&row.Article,
			&row.Caption,
			&row.TypeID,
			&row.Length,
			&row.Width,
			&row.Thickness,
			&row.Weight,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// FetchTypeIDs - comment method.
func (re *LaminatePostgres) FetchTypeIDs(ctx context.Context) ([]uint64, error) {
	return re.repoTypeIDs.Fetch(ctx, mrenum.ItemStatusEnabled)
}

// FetchThicknesses - comment method.
func (re *LaminatePostgres) FetchThicknesses(ctx context.Context) ([]measure.Meter, error) {
	return re.repoThicknesses.Fetch(ctx, mrenum.ItemStatusEnabled)
}

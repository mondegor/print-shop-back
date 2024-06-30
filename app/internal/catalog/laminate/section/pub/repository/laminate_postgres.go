package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// LaminatePostgres - comment struct.
	// LaminatePostgres - comment struct.
	LaminatePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewLaminatePostgres - создаёт объект LaminatePostgres.
func NewLaminatePostgres(client mrstorage.DBConnManager) *LaminatePostgres {
	return &LaminatePostgres{
		client: client,
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
			laminate_weight
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameLaminates + `
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
func (re *LaminatePostgres) FetchTypeIDs(ctx context.Context) ([]mrtype.KeyInt32, error) {
	sql := `
        SELECT
			type_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameLaminates + `
        WHERE
            laminate_status = $1 AND deleted_at IS NULL
        GROUP BY
            type_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]mrtype.KeyInt32, 0)

	for cursor.Next() {
		var typeID mrtype.KeyInt32

		err = cursor.Scan(
			&typeID,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, typeID)
	}

	return rows, cursor.Err()
}

// FetchThicknesses - comment method.
func (re *LaminatePostgres) FetchThicknesses(ctx context.Context) ([]measure.Micrometer, error) {
	sql := `
        SELECT
			laminate_thickness
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameLaminates + `
        WHERE
            laminate_status = $1 AND deleted_at IS NULL
        GROUP BY
            laminate_thickness ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]measure.Micrometer, 0)

	for cursor.Next() {
		var thickness measure.Micrometer

		err = cursor.Scan(
			&thickness,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, thickness)
	}

	return rows, cursor.Err()
}

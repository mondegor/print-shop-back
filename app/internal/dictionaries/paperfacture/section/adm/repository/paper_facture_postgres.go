package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/shared/repository"
)

type (
	// PaperFacturePostgres - comment struct.
	PaperFacturePostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewPaperFacturePostgres - создаёт объект PaperFacturePostgres.
func NewPaperFacturePostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *PaperFacturePostgres) NewSelectParams(params entity.PaperFactureParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLike("UPPER(facture_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("facture_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("facture_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *PaperFacturePostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.PaperFacture, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
        SELECT
            facture_id,
            tag_version,
            facture_caption as caption,
            facture_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + params.OrderBy.String() + params.Limit.String() + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PaperFacture, 0)

	for cursor.Next() {
		var row entity.PaperFacture

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Caption,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// FetchTotal - comment method.
func (re *PaperFacturePostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        WHERE
            ` + whereStr + `;`

	var totalRow int64

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		whereArgs...,
	).Scan(
		&totalRow,
	)

	return totalRow, err
}

// FetchOne - comment method.
func (re *PaperFacturePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperFacture, error) {
	sql := `
        SELECT
            tag_version,
            facture_caption,
            facture_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        WHERE
            facture_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.PaperFacture{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository.PaperFactureFetchStatusPostgres(ctx, re.client, rowID)
}

// Insert - comment method.
func (re *PaperFacturePostgres) Insert(ctx context.Context, row entity.PaperFacture) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
            (
                facture_caption,
                facture_status
            )
        VALUES
            ($1, $2)
        RETURNING
            facture_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *PaperFacturePostgres) Update(ctx context.Context, row entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_caption = $3
        WHERE
            facture_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *PaperFacturePostgres) UpdateStatus(ctx context.Context, row entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_status = $3
        WHERE
            facture_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// Delete - comment method.
func (re *PaperFacturePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            facture_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

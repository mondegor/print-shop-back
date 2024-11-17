package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
)

type (
	// PaperColorPostgres - comment struct.
	PaperColorPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewPaperColorPostgres - создаёт объект PaperColorPostgres.
func NewPaperColorPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *PaperColorPostgres {
	return &PaperColorPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNamePaperColors,
			"color_id",
			module.DBFieldTagVersion,
			"color_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNamePaperColors,
			"color_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNamePaperColors,
		),
	}
}

// FetchWithTotal - comment method.
func (re *PaperColorPostgres) FetchWithTotal(ctx context.Context, params entity.PaperColorParams) (rows []entity.PaperColor, countRows uint64, err error) {
	condition := re.sqlBuilder.Condition().Build(re.fetchCondition(params.Filter))

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().Build(re.fetchOrderBy(params.Sorter))
	limit := re.sqlBuilder.Limit().Build(params.Pager.Index, params.Pager.Size)

	rows, err = re.fetch(ctx, condition, orderBy, limit, params.Pager.Size)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *PaperColorPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.PaperColor, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
        SELECT
            color_id,
            tag_version,
            color_caption as caption,
            color_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBTableNamePaperColors + `
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + orderBy.String() + limit.String() + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PaperColor, 0, maxRows)

	for cursor.Next() {
		var row entity.PaperColor

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

// FetchOne - comment method.
func (re *PaperColorPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.PaperColor, error) {
	sql := `
        SELECT
            tag_version,
            color_caption,
            color_status,
            created_at,
			updated_at
        FROM
            ` + module.DBTableNamePaperColors + `
        WHERE
            color_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.PaperColor{ID: rowID}

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
func (re *PaperColorPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *PaperColorPostgres) Insert(ctx context.Context, row entity.PaperColor) (rowID uint64, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNamePaperColors + `
            (
                color_caption,
                color_status
            )
        VALUES
            ($1, $2)
        RETURNING
            color_id;`

	err = re.client.Conn(ctx).QueryRow(
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
func (re *PaperColorPostgres) Update(ctx context.Context, row entity.PaperColor) (tagVersion uint32, err error) {
	sql := `
        UPDATE
            ` + module.DBTableNamePaperColors + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            color_caption = $3
        WHERE
            color_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	err = re.client.Conn(ctx).QueryRow(
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
func (re *PaperColorPostgres) UpdateStatus(ctx context.Context, row entity.PaperColor) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *PaperColorPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *PaperColorPostgres) fetchCondition(filter entity.PaperColorListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(color_caption)", strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("color_status", filter.Statuses),
			)
		},
	)
}

func (re *PaperColorPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("color_id", mrenum.SortDirectionASC),
			)
		},
	)
}

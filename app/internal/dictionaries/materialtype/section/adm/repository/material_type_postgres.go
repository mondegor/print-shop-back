package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtype/sortdirection"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, itemstatus.Enum]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[int]
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, itemstatus.Enum](
			client,
			module.DBTableNameMaterialTypes,
			"type_id",
			module.DBFieldTagVersion,
			"type_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameMaterialTypes,
			"type_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[int](
			client,
			module.DBTableNameMaterialTypes,
		),
	}
}

// FetchWithTotal - comment method.
func (re *MaterialTypePostgres) FetchWithTotal(
	ctx context.Context,
	params entity.MaterialTypeParams,
) (rows []entity.MaterialType, countRows int, err error) {
	condition := re.sqlBuilder.Condition().BuildFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(type_caption)", strings.ToUpper(params.Filter.SearchText)),
				c.FilterAnyOf("type_status", params.Filter.Statuses),
			)
		},
	)

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().BuildFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(params.Sorter.Column, params.Sorter.Direction),
				o.Field("type_id", sortdirection.ASC),
			)
		},
	)
	limit := re.sqlBuilder.Limit().Build(params.Pager.Index, params.Pager.Size)

	rows, err = re.fetch(ctx, condition, orderBy, limit, params.Pager.Size)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *MaterialTypePostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows int,
) ([]entity.MaterialType, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
        SELECT
            type_id,
            tag_version,
            type_caption as caption,
            type_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBTableNameMaterialTypes + `
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + mrstorage.ToSQL(orderBy) + mrstorage.ToSQL(limit) + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.MaterialType, 0, maxRows)

	for cursor.Next() {
		var row entity.MaterialType

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
func (re *MaterialTypePostgres) FetchOne(ctx context.Context, rowID uint64) (entity.MaterialType, error) {
	sql := `
        SELECT
            tag_version,
            type_caption,
            type_status,
            created_at,
			updated_at
        FROM
            ` + module.DBTableNameMaterialTypes + `
        WHERE
            type_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.MaterialType{ID: rowID}

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
// result: itemstatus.Enum - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *MaterialTypePostgres) FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *MaterialTypePostgres) Insert(ctx context.Context, row entity.MaterialType) (rowID uint64, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNameMaterialTypes + `
            (
                type_caption,
                type_status
            )
        VALUES
            ($1, $2)
        RETURNING
            type_id;`

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
func (re *MaterialTypePostgres) Update(ctx context.Context, row entity.MaterialType) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// UpdateStatus - comment method.
func (re *MaterialTypePostgres) UpdateStatus(ctx context.Context, row entity.MaterialType) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *MaterialTypePostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

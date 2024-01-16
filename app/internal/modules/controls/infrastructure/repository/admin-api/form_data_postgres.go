package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormDataPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewFormDataPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *FormDataPostgres {
	return &FormDataPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *FormDataPostgres) NewFetchParams(params entity.FormDataParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("form_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(param_name)", "UPPER(form_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("form_detailing", params.Filter.Detailing),
				w.FilterAnyOf("form_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("form_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *FormDataPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.FormData, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            form_id,
            tag_version,
            datetime_created as createdAt,
			datetime_updated as updatedAt,
            param_name as paramName,
            form_caption as caption,
            form_detailing,
            form_status
        FROM
            ` + module.UnitFormDataDBSchema + `.forms
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + params.OrderBy.String() + params.Pager.String() + `;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		whereArgs...,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.FormData, 0)

	for cursor.Next() {
		var row entity.FormData

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.ParamName,
			&row.Caption,
			&row.Detailing,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *FormDataPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.UnitFormDataDBSchema + `.forms
        WHERE
            ` + whereStr + `;`

	var totalRow int64

	err := re.client.QueryRow(
		ctx,
		sql,
		whereArgs...,
	).Scan(
		&totalRow,
	)

	return totalRow, err
}

func (re *FormDataPostgres) LoadOne(ctx context.Context, row *entity.FormData) error {
	sql := `
        SELECT
            tag_version,
            datetime_created,
			datetime_updated,
            param_name,
            form_caption,
            form_detailing,
            form_body_compiled,
            form_status
        FROM
            ` + module.UnitFormDataDBSchema + `.forms
        WHERE
            form_id = $1 AND form_status <> $2
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.ParamName,
		&row.Caption,
		&row.Detailing,
		&row.Body,
		&row.Status,
	)
}

func (re *FormDataPostgres) FetchIdByName(ctx context.Context, paramName string) (mrtype.KeyInt32, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements
        WHERE
            param_name = $1
        LIMIT 1;`

	var id mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		paramName,
	).Scan(
		&id,
	)

	return id, err
}

func (re *FormDataPostgres) FetchStatus(ctx context.Context, row *entity.FormData) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            form_status
        FROM
            ` + module.UnitFormDataDBSchema + `.forms
        WHERE
            form_id = $1 AND form_status <> $2
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *FormDataPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        SELECT
            1
        FROM
            ` + module.UnitFormDataDBSchema + `.forms
        WHERE
            form_id = $1 AND form_status <> $2
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	).Scan(
		&id,
	)
}

func (re *FormDataPostgres) Insert(ctx context.Context, row *entity.FormData) error {
	sql := `
        INSERT INTO ` + module.UnitFormDataDBSchema + `.forms
            (
                param_name,
                form_caption,
                form_detailing,
                form_body_compiled,
                form_status
            )
        VALUES
            ($1, $2, $3, '[]', $4)
        RETURNING
            form_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ParamName,
		row.Caption,
		row.Detailing,
		row.Status,
	).Scan(
		&row.ID,
	)

	return err
}

func (re *FormDataPostgres) Update(ctx context.Context, row *entity.FormData) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
	}

	setStr, setArgs := set.Param(len(args) + 1).ToSql()

	sql := `
        UPDATE
            ` + module.UnitFormDataDBSchema + `.forms
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
            ` + setStr + `
        WHERE
            form_id = $1 AND tag_version = $2 AND form_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *FormDataPostgres) UpdateStatus(ctx context.Context, row *entity.FormData) (int32, error) {
	sql := `
        UPDATE
            ` + module.UnitFormDataDBSchema + `.forms
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
            form_status = $4
        WHERE
            form_id = $1 AND tag_version = $2 AND form_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *FormDataPostgres) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.UnitFormDataDBSchema + `.forms
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
            param_name = NULL,
            form_status = $2
        WHERE
            form_id = $1 AND form_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	)
}

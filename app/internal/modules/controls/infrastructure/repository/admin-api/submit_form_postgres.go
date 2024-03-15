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
	SubmitFormPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewSubmitFormPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *SubmitFormPostgres {
	return &SubmitFormPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *SubmitFormPostgres) NewFetchParams(params entity.SubmitFormParams) mrstorage.SqlSelectParams {
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

func (re *SubmitFormPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.SubmitForm, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            form_id,
            tag_version,
            created_at as createdAt,
			updated_at as updatedAt,
            param_name as paramName,
            form_caption as caption,
            form_detailing,
            form_status
        FROM
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
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

	rows := make([]entity.SubmitForm, 0)

	for cursor.Next() {
		var row entity.SubmitForm

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

func (re *SubmitFormPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
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

func (re *SubmitFormPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.SubmitForm, error) {
	sql := `
        SELECT
            tag_version,
            created_at,
			updated_at,
            param_name,
            form_caption,
            form_detailing,
            form_body_compiled,
            form_status
        FROM
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
        WHERE
            form_id = $1 AND form_status <> $2
        LIMIT 1;`

	row := entity.SubmitForm{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
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

	return row, err
}

func (re *SubmitFormPostgres) FetchIdByName(ctx context.Context, paramName string) (mrtype.KeyInt32, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.UnitFormElementsDBSchema + `.submit_form_elements
        WHERE
            param_name = $1
        LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		paramName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

func (re *SubmitFormPostgres) FetchStatus(ctx context.Context, row entity.SubmitForm) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            form_status
        FROM
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
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
func (re *SubmitFormPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
        WHERE
            form_id = $1 AND form_status <> $2
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}

func (re *SubmitFormPostgres) Insert(ctx context.Context, row entity.SubmitForm) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.UnitSubmitFormDBSchema + `.submit_forms
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

	return row.ID, err
}

func (re *SubmitFormPostgres) Update(ctx context.Context, row entity.SubmitForm) (int32, error) {
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
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
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

func (re *SubmitFormPostgres) UpdateStatus(ctx context.Context, row entity.SubmitForm) (int32, error) {
	sql := `
        UPDATE
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
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

func (re *SubmitFormPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.UnitSubmitFormDBSchema + `.submit_forms
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            param_name = NULL,
            form_status = $2
        WHERE
            form_id = $1 AND form_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}

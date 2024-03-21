package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls/submit-form"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
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

func (re *SubmitFormPostgres) NewSelectParams(params entity.SubmitFormParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("form_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(param_name)", "UPPER(form_caption)"}, strings.ToUpper(params.Filter.SearchText)),
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
			rewrite_name as rewriteName,
            param_name as paramName,
            form_caption as caption,
            form_detailing,
            form_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.submit_forms
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
			&row.RewriteName,
			&row.ParamName,
			&row.Caption,
			&row.Detailing,
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

func (re *SubmitFormPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.submit_forms
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

func (re *SubmitFormPostgres) FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error) {
	sql := `
        SELECT
            tag_version,
			rewrite_name,
            param_name,
            form_caption,
            form_detailing,
            form_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.submit_forms
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
		&row.RewriteName,
		&row.ParamName,
		&row.Caption,
		&row.Detailing,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (re *SubmitFormPostgres) FetchIdByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.DBSchema + `.submit_forms
        WHERE
            rewrite_name = $1
        LIMIT 1;`

	var rowID uuid.UUID

	err := re.client.QueryRow(
		ctx,
		sql,
		rewriteName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

func (re *SubmitFormPostgres) FetchIdByParamName(ctx context.Context, paramName string) (uuid.UUID, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.DBSchema + `.submit_forms
        WHERE
            param_name = $1
        LIMIT 1;`

	var rowID uuid.UUID

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
            ` + module.DBSchema + `.submit_forms
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
func (re *SubmitFormPostgres) IsExists(ctx context.Context, rowID uuid.UUID) error {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.DBSchema + `.submit_forms
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

func (re *SubmitFormPostgres) Insert(ctx context.Context, row entity.SubmitForm) (uuid.UUID, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.submit_forms
            (
				form_id,
				rewrite_name,
                param_name,
                form_caption,
                form_detailing,
                form_status
            )
        VALUES
            (gen_random_uuid(), $1, $2, $3, $4, $5)
        RETURNING
            form_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.RewriteName,
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
            ` + module.DBSchema + `.submit_forms
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
            ` + module.DBSchema + `.submit_forms
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

func (re *SubmitFormPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.submit_forms
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
			rewrite_name = NULL,
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

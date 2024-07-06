package repository

import (
	"context"
	"strings"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// SubmitFormPostgres - comment struct.
	SubmitFormPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewSubmitFormPostgres - создаёт объект SubmitFormPostgres.
func NewSubmitFormPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect, sqlUpdate mrstorage.SQLBuilderUpdate) *SubmitFormPostgres {
	return &SubmitFormPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewSelectParams - comment method.
func (re *SubmitFormPostgres) NewSelectParams(params entity.SubmitFormParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(param_name)", "UPPER(form_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("form_detailing", params.Filter.Detailing),
				w.FilterAnyOf("form_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("form_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *SubmitFormPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.SubmitForm, error) {
	whereStr, whereArgs := params.Where.ToSQL()

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
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
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

// FetchTotal - comment method.
func (re *SubmitFormPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
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
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        WHERE
            form_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.SubmitForm{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
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

// FetchIDByRewriteName - comment method.
func (re *SubmitFormPostgres) FetchIDByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        WHERE
            rewrite_name = $1 AND deleted_at IS NULL
        LIMIT 1;`

	var rowID uuid.UUID

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rewriteName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// FetchIDByParamName - comment method.
func (re *SubmitFormPostgres) FetchIDByParamName(ctx context.Context, paramName string) (uuid.UUID, error) {
	sql := `
        SELECT
            form_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        WHERE
            param_name = $1 AND deleted_at IS NULL
        LIMIT 1;`

	var rowID uuid.UUID

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		paramName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *SubmitFormPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            form_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        WHERE
            form_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}

// Insert - comment method.
func (re *SubmitFormPostgres) Insert(ctx context.Context, row entity.SubmitForm) (uuid.UUID, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
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

	err := re.client.Conn(ctx).QueryRow(
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

// Update - comment method.
func (re *SubmitFormPostgres) Update(ctx context.Context, row entity.SubmitForm) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.WithParam(len(args) + 1).ToSQL()

	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            form_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *SubmitFormPostgres) UpdateStatus(ctx context.Context, row entity.SubmitForm) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            form_status = $3
        WHERE
            form_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *SubmitFormPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameSubmitForms + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            form_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

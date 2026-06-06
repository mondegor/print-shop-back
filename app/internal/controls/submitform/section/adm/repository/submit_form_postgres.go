package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrpostgres/db"
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-sysmess/mrstorage/mrsql"
	"github.com/mondegor/go-sysmess/mrtype/sortdirection"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/controls/submitform/module"
	"print-shop-back/internal/controls/submitform/section/adm/entity"
)

type (
	// SubmitFormPostgres - comment struct.
	SubmitFormPostgres struct {
		client              mrstorage.DBConnManager
		sqlBuilder          mrstorage.SQLBuilder
		repoIDByRewriteName db.FieldFetcher[string, uuid.UUID]
		repoIDByParamName   db.FieldFetcher[string, uuid.UUID]
		repoStatus          db.FieldWithVersionUpdater[uuid.UUID, uint32, workflow.ItemStatus]
		repoSoftDeleter     db.RowSoftDeleter[uuid.UUID]
		repoTotalRows       db.TotalRowsFetcher[int]
	}
)

// NewSubmitFormPostgres - создаёт объект SubmitFormPostgres.
func NewSubmitFormPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *SubmitFormPostgres {
	return &SubmitFormPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoIDByRewriteName: db.NewFieldFetcher[string, uuid.UUID](
			client,
			module.DBTableNameSubmitForms,
			"rewrite_name",
			"form_id",
			module.DBFieldDeletedAt,
		),
		repoIDByParamName: db.NewFieldFetcher[string, uuid.UUID](
			client,
			module.DBTableNameSubmitForms,
			"param_name",
			"form_id",
			module.DBFieldDeletedAt,
		),
		repoStatus: db.NewFieldWithVersionUpdater[uuid.UUID, uint32, workflow.ItemStatus](
			client,
			module.DBTableNameSubmitForms,
			"form_id",
			module.DBFieldTagVersion,
			"form_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uuid.UUID](
			client,
			module.DBTableNameSubmitForms,
			"form_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[int](
			client,
			module.DBTableNameSubmitForms,
		),
	}
}

// FetchWithTotal - comment method.
func (re *SubmitFormPostgres) FetchWithTotal(ctx context.Context, params entity.SubmitFormParams) (rows []entity.SubmitForm, countRows int, err error) {
	condition := re.sqlBuilder.Condition().BuildFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(param_name)", "UPPER(form_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				c.FilterAnyOf("form_detailing", params.Filter.Detailing),
				c.FilterAnyOf("form_status", params.Filter.Statuses),
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
				o.Field("form_id", sortdirection.ASC),
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
func (re *SubmitFormPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows int,
) ([]entity.SubmitForm, error) {
	whereStr, whereArgs := condition.ToSQL()

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
            ` + module.DBTableNameSubmitForms + `
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

	rows := make([]entity.SubmitForm, 0, maxRows)

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
            ` + module.DBTableNameSubmitForms + `
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
func (re *SubmitFormPostgres) FetchIDByRewriteName(ctx context.Context, rewriteName string) (rowID uuid.UUID, err error) {
	return re.repoIDByRewriteName.Fetch(ctx, rewriteName)
}

// FetchIDByParamName - comment method.
func (re *SubmitFormPostgres) FetchIDByParamName(ctx context.Context, paramName string) (rowID uuid.UUID, err error) {
	return re.repoIDByParamName.Fetch(ctx, paramName)
}

// FetchStatus - comment method.
// result: workflow.ItemStatus - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *SubmitFormPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (workflow.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *SubmitFormPostgres) Insert(ctx context.Context, row entity.SubmitForm) (rowID uuid.UUID, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNameSubmitForms + `
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

	err = re.client.Conn(ctx).QueryRow(
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
func (re *SubmitFormPostgres) Update(ctx context.Context, row entity.SubmitForm) (tagVersion uint32, err error) {
	set, err := re.sqlBuilder.Set().BuildEntity(row)
	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.WithStartArg(len(args) + 1).ToSQL()

	sql := `
        UPDATE
            ` + module.DBTableNameSubmitForms + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            form_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

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
func (re *SubmitFormPostgres) UpdateStatus(ctx context.Context, row entity.SubmitForm) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *SubmitFormPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

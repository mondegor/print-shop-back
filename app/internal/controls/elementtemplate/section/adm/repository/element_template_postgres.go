package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-core/mrpostgres/db"
	"github.com/mondegor/go-core/mrstorage"
	"github.com/mondegor/go-core/mrstorage/mrsql"
	"github.com/mondegor/go-core/mrtype/sortdirection"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/controls/elementtemplate/module"
	"print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplatePostgres - comment struct.
	ElementTemplatePostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, workflow.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[int]
	}
)

// NewElementTemplatePostgres - создаёт объект ElementTemplatePostgres.
func NewElementTemplatePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *ElementTemplatePostgres {
	return &ElementTemplatePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, workflow.ItemStatus](
			client,
			module.DBTableNameElementTemplates,
			"template_id",
			module.DBFieldTagVersion,
			"template_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameElementTemplates,
			"template_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[int](
			client,
			module.DBTableNameElementTemplates,
		),
	}
}

// FetchWithTotal - comment method.
func (re *ElementTemplatePostgres) FetchWithTotal(
	ctx context.Context,
	params entity.ElementTemplateParams,
) (rows []entity.ElementTemplate, countRows int, err error) {
	condition := re.sqlBuilder.Condition().BuildFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLikeFields([]string{"UPPER(param_name)", "UPPER(template_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				c.FilterAnyOf("element_detailing", params.Filter.Detailing),
				c.FilterAnyOf("template_status", params.Filter.Statuses),
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
				o.Field("template_id", sortdirection.ASC),
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
func (re *ElementTemplatePostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows int,
) ([]entity.ElementTemplate, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
        SELECT
            template_id,
            tag_version,
            param_name as paramName,
            template_caption as caption,
            element_type,
            element_detailing,
            template_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBTableNameElementTemplates + `
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

	rows := make([]entity.ElementTemplate, 0, maxRows)

	for cursor.Next() {
		var row entity.ElementTemplate

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.ParamName,
			&row.Caption,
			&row.Type,
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
func (re *ElementTemplatePostgres) FetchOne(ctx context.Context, rowID uint64) (entity.ElementTemplate, error) {
	sql := `
        SELECT
            tag_version,
            param_name,
            template_caption,
            element_type,
            element_detailing,
            element_body,
            template_status,
            created_at,
			updated_at
        FROM
            ` + module.DBTableNameElementTemplates + `
        WHERE
            template_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.ElementTemplate{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.ParamName,
		&row.Caption,
		&row.Type,
		&row.Detailing,
		&row.Body,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: workflow.ItemStatus - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *ElementTemplatePostgres) FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *ElementTemplatePostgres) Insert(ctx context.Context, row entity.ElementTemplate) (rowID uint64, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNameElementTemplates + `
            (
                param_name,
                template_caption,
                element_type,
                element_detailing,
                element_body,
                template_status
            )
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING
            template_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ParamName,
		row.Caption,
		row.Type,
		row.Detailing,
		row.Body,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *ElementTemplatePostgres) Update(ctx context.Context, row entity.ElementTemplate) (tagVersion uint32, err error) {
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
            ` + module.DBTableNameElementTemplates + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            template_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *ElementTemplatePostgres) UpdateStatus(ctx context.Context, row entity.ElementTemplate) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *ElementTemplatePostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

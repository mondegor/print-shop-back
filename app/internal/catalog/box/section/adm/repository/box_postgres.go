package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
)

type (
	// BoxPostgres - comment struct.
	BoxPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoIDByArticle db.FieldFetcher[string, uint64]
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewBoxPostgres - создаёт объект BoxPostgres.
func NewBoxPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *BoxPostgres {
	return &BoxPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoIDByArticle: db.NewFieldFetcher[string, uint64](
			client,
			module.DBTableNameBoxes,
			"box_article",
			"box_id",
			module.DBFieldDeletedAt,
		),
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNameBoxes,
			"box_id",
			module.DBFieldTagVersion,
			"box_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameBoxes,
			"box_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameBoxes,
		),
	}
}

// FetchWithTotal - comment method.
func (re *BoxPostgres) FetchWithTotal(ctx context.Context, params entity.BoxParams) (rows []entity.Box, countRows uint64, err error) {
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
func (re *BoxPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Box, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
		SELECT
			box_id,
			tag_version,
			box_article as article,
			box_caption as caption,
			box_length as length,
			box_width as width,
			box_height as height,
			box_thickness as thickness,
			box_weight as weight,
			box_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBTableNameBoxes + `
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

	rows := make([]entity.Box, 0, maxRows)

	for cursor.Next() {
		var row entity.Box

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Article,
			&row.Caption,
			&row.Length,
			&row.Width,
			&row.Height,
			&row.Thickness,
			&row.Weight,
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
func (re *BoxPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.Box, error) {
	sql := `
		SELECT
			tag_version,
			box_article,
			box_caption,
			box_length,
			box_width,
			box_height,
			box_thickness,
			box_weight,
			box_status,
			created_at,
			updated_at
		FROM
			` + module.DBTableNameBoxes + `
		WHERE
			box_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Box{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Article,
		&row.Caption,
		&row.Length,
		&row.Width,
		&row.Height,
		&row.Thickness,
		&row.Weight,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchIDByArticle - comment method.
func (re *BoxPostgres) FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error) {
	return re.repoIDByArticle.Fetch(ctx, article)
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *BoxPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *BoxPostgres) Insert(ctx context.Context, row entity.Box) (rowID uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameBoxes + `
			(
				box_article,
				box_caption,
				box_length,
				box_width,
				box_height,
				box_thickness,
				box_weight,
				box_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			box_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.Length,
		row.Width,
		row.Height,
		row.Thickness,
		row.Weight,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *BoxPostgres) Update(ctx context.Context, row entity.Box) (tagVersion uint32, err error) {
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
			` + module.DBTableNameBoxes + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			box_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *BoxPostgres) UpdateStatus(ctx context.Context, row entity.Box) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *BoxPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *BoxPostgres) fetchCondition(filter entity.BoxListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLikeFields([]string{"UPPER(box_article)", "UPPER(box_caption)"}, strings.ToUpper(filter.SearchText)),
				c.FilterRangeFloat64("box_length", mrtype.RangeFloat64(filter.Length), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("box_width", mrtype.RangeFloat64(filter.Width), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("box_height", mrtype.RangeFloat64(filter.Height), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("box_weight", mrtype.RangeFloat64(filter.Weight), 0, mrlib.EqualityThresholdE9),
				c.FilterAnyOf("box_status", filter.Statuses),
			)
		},
	)
}

func (re *BoxPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("box_id", mrenum.SortDirectionASC),
			)
		},
	)
}

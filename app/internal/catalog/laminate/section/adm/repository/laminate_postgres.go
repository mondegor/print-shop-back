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

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
)

type (
	// LaminatePostgres - comment struct.
	LaminatePostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoIDByArticle db.FieldFetcher[string, uint64]
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewLaminatePostgres - создаёт объект LaminatePostgres.
func NewLaminatePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *LaminatePostgres {
	return &LaminatePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNameLaminates,
			"laminate_id",
			module.DBFieldTagVersion,
			"laminate_status",
			module.DBFieldDeletedAt,
		),
		repoIDByArticle: db.NewFieldFetcher[string, uint64](
			client,
			module.DBTableNameLaminates,
			"laminate_article",
			"laminate_id",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameLaminates,
			"laminate_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameLaminates,
		),
	}
}

// FetchWithTotal - comment method.
func (re *LaminatePostgres) FetchWithTotal(ctx context.Context, params entity.LaminateParams) (rows []entity.Laminate, countRows uint64, err error) {
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
func (re *LaminatePostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Laminate, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
		SELECT
			laminate_id,
			tag_version,
			laminate_article as article,
			laminate_caption as caption,
			type_id,
			laminate_length as length,
			laminate_width as width,
			laminate_thickness,
			laminate_weight_m2 as weightM2,
			laminate_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBTableNameLaminates + `
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

	rows := make([]entity.Laminate, 0, maxRows)

	for cursor.Next() {
		var row entity.Laminate

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Article,
			&row.Caption,
			&row.TypeID,
			&row.Length,
			&row.Width,
			&row.Thickness,
			&row.WeightM2,
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
func (re *LaminatePostgres) FetchOne(ctx context.Context, rowID uint64) (entity.Laminate, error) {
	sql := `
		SELECT
			tag_version,
			laminate_article,
			laminate_caption,
			type_id,
			laminate_length,
			laminate_width,
			laminate_thickness,
			laminate_weight_m2,
			laminate_status,
			created_at,
			updated_at
		FROM
			` + module.DBTableNameLaminates + `
		WHERE
			laminate_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Laminate{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Article,
		&row.Caption,
		&row.TypeID,
		&row.Length,
		&row.Width,
		&row.Thickness,
		&row.WeightM2,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchIDByArticle - comment method.
func (re *LaminatePostgres) FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error) {
	return re.repoIDByArticle.Fetch(ctx, article)
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *LaminatePostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *LaminatePostgres) Insert(ctx context.Context, row entity.Laminate) (rowID uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameLaminates + `
			(
				laminate_article,
				laminate_caption,
				type_id,
				laminate_length,
				laminate_width,
				laminate_thickness,
				laminate_weight_m2,
				laminate_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			laminate_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.TypeID,
		row.Length,
		row.Width,
		row.Thickness,
		row.WeightM2,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *LaminatePostgres) Update(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error) {
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
			` + module.DBTableNameLaminates + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			laminate_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *LaminatePostgres) UpdateStatus(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *LaminatePostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *LaminatePostgres) fetchCondition(filter entity.LaminateListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLikeFields([]string{"UPPER(laminate_article)", "UPPER(laminate_caption)"}, strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("type_id", filter.TypeIDs),
				c.FilterRangeFloat64("laminate_length", mrtype.RangeFloat64(filter.Length), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("laminate_width", mrtype.RangeFloat64(filter.Width), 0, mrlib.EqualityThresholdE9),
				c.FilterAnyOf("laminate_status", filter.Statuses),
			)
		},
	)
}

func (re *LaminatePostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("laminate_id", mrenum.SortDirectionASC),
			)
		},
	)
}

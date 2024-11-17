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

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
)

type (
	// PaperPostgres - comment struct.
	PaperPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoIDByArticle db.FieldFetcher[string, uint64]
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewPaperPostgres - создаёт объект PaperPostgres.
func NewPaperPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *PaperPostgres {
	return &PaperPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoIDByArticle: db.NewFieldFetcher[string, uint64](
			client,
			module.DBTableNamePapers,
			"paper_article",
			"paper_id",
			module.DBFieldDeletedAt,
		),
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNamePapers,
			"paper_id",
			module.DBFieldTagVersion,
			"paper_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNamePapers,
			"paper_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNamePapers,
		),
	}
}

// FetchWithTotal - comment method.
func (re *PaperPostgres) FetchWithTotal(ctx context.Context, params entity.PaperParams) (rows []entity.Paper, countRows uint64, err error) {
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
func (re *PaperPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Paper, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
		SELECT
			paper_id,
			tag_version,
			paper_article as article,
			paper_caption as caption,
			type_id,
			color_id,
			facture_id,
			paper_width as width,
			paper_height as height,
			paper_thickness,
			paper_density as density,
			paper_sides,
			paper_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBTableNamePapers + `
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

	rows := make([]entity.Paper, 0, maxRows)

	for cursor.Next() {
		var row entity.Paper

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Article,
			&row.Caption,
			&row.TypeID,
			&row.ColorID,
			&row.FactureID,
			&row.Width,
			&row.Height,
			&row.Thickness,
			&row.Density,
			&row.Sides,
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
func (re *PaperPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.Paper, error) {
	sql := `
		SELECT
			tag_version,
			paper_article,
			paper_caption,
			type_id,
			color_id,
			facture_id,
			paper_width,
			paper_height,
			paper_thickness,
			paper_density,
			paper_sides,
			paper_status,
			created_at,
			updated_at
		FROM
			` + module.DBTableNamePapers + `
		WHERE
			paper_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Paper{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Article,
		&row.Caption,
		&row.TypeID,
		&row.ColorID,
		&row.FactureID,
		&row.Width,
		&row.Height,
		&row.Thickness,
		&row.Density,
		&row.Sides,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchIDByArticle - comment method.
func (re *PaperPostgres) FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error) {
	return re.repoIDByArticle.Fetch(ctx, article)
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *PaperPostgres) Insert(ctx context.Context, row entity.Paper) (rowID uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNamePapers + `
			(
				paper_article,
				paper_caption,
				type_id,
				color_id,
				facture_id,
				paper_width,
				paper_height,
				paper_thickness,
				paper_density,
				paper_sides,
				paper_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING
			paper_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.TypeID,
		row.ColorID,
		row.FactureID,
		row.Width,
		row.Height,
		row.Thickness,
		row.Density,
		row.Sides,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *PaperPostgres) Update(ctx context.Context, row entity.Paper) (tagVersion uint32, err error) {
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
			` + module.DBTableNamePapers + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			paper_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *PaperPostgres) UpdateStatus(ctx context.Context, row entity.Paper) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *PaperPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *PaperPostgres) fetchCondition(filter entity.PaperListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLikeFields([]string{"UPPER(paper_article)", "UPPER(paper_caption)"}, strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("type_id", filter.TypeIDs),
				c.FilterAnyOf("color_id", filter.ColorIDs),
				c.FilterAnyOf("facture_id", filter.FactureIDs),
				c.FilterRangeFloat64("paper_width", mrtype.RangeFloat64(filter.Width), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("paper_height", mrtype.RangeFloat64(filter.Height), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("paper_density", mrtype.RangeFloat64(filter.Density), 0, mrlib.EqualityThresholdE9),
				c.FilterAnyOf("paper_status", filter.Statuses),
			)
		},
	)
}

func (re *PaperPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("paper_id", mrenum.SortDirectionASC),
			)
		},
	)
}

package repository

import (
	"context"
	"strings"

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
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewLaminatePostgres - создаёт объект LaminatePostgres.
func NewLaminatePostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect, sqlUpdate mrstorage.SQLBuilderUpdate) *LaminatePostgres {
	return &LaminatePostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewSelectParams - comment method.
func (re *LaminatePostgres) NewSelectParams(params entity.LaminateParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLikeFields([]string{"UPPER(laminate_article)", "UPPER(laminate_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("type_id", params.Filter.TypeIDs),
				w.FilterRangeFloat64("laminate_length", mrtype.RangeFloat64(params.Filter.Length), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("laminate_width", mrtype.RangeFloat64(params.Filter.Width), 0, mrlib.EqualityThresholdE9),
				w.FilterAnyOf("laminate_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("laminate_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *LaminatePostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Laminate, error) {
	whereStr, whereArgs := params.Where.ToSQL()

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
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
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

	rows := make([]entity.Laminate, 0)

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

// FetchTotal - comment method.
func (re *LaminatePostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
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
func (re *LaminatePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Laminate, error) {
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
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
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
func (re *LaminatePostgres) FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			laminate_id
		FROM
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
		WHERE
			laminate_article = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *LaminatePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			laminate_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
		WHERE
			laminate_id = $1 AND deleted_at IS NULL
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
func (re *LaminatePostgres) Insert(ctx context.Context, row entity.Laminate) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameLaminates + `
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
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			laminate_id;`

	err := re.client.Conn(ctx).QueryRow(
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
func (re *LaminatePostgres) Update(ctx context.Context, row entity.Laminate) (int32, error) {
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
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			laminate_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *LaminatePostgres) UpdateStatus(ctx context.Context, row entity.Laminate) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			laminate_status = $3
		WHERE
			laminate_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *LaminatePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameLaminates + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			laminate_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

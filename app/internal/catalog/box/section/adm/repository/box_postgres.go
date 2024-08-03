package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// BoxPostgres - comment struct.
	BoxPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewBoxPostgres - создаёт объект BoxPostgres.
func NewBoxPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect, sqlUpdate mrstorage.SQLBuilderUpdate) *BoxPostgres {
	return &BoxPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewSelectParams - comment method.
func (re *BoxPostgres) NewSelectParams(params entity.BoxParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLikeFields([]string{"UPPER(box_article)", "UPPER(box_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterRangeFloat64("box_length", mrtype.RangeFloat64(params.Filter.Length), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("box_width", mrtype.RangeFloat64(params.Filter.Width), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("box_height", mrtype.RangeFloat64(params.Filter.Height), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("box_weight", mrtype.RangeFloat64(params.Filter.Weight), 0, mrlib.EqualityThresholdE9),
				w.FilterAnyOf("box_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("box_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *BoxPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Box, error) {
	whereStr, whereArgs := params.Where.ToSQL()

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
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
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

	rows := make([]entity.Box, 0)

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

// FetchTotal - comment method.
func (re *BoxPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
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
func (re *BoxPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Box, error) {
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
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
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
func (re *BoxPostgres) FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			box_id
		FROM
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
		WHERE
			box_article = $1 AND deleted_at IS NULL
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
func (re *BoxPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			box_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
		WHERE
			box_id = $1 AND deleted_at IS NULL
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
func (re *BoxPostgres) Insert(ctx context.Context, row entity.Box) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameBoxes + `
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

	err := re.client.Conn(ctx).QueryRow(
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
func (re *BoxPostgres) Update(ctx context.Context, row entity.Box) (int32, error) {
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
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			box_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *BoxPostgres) UpdateStatus(ctx context.Context, row entity.Box) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			box_status = $3
		WHERE
			box_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *BoxPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameBoxes + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			box_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

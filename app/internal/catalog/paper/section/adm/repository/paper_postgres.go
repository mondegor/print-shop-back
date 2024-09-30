package repository

import (
	"context"
	"strings"

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
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewPaperPostgres - создаёт объект PaperPostgres.
func NewPaperPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect, sqlUpdate mrstorage.SQLBuilderUpdate) *PaperPostgres {
	return &PaperPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewSelectParams - comment method.
func (re *PaperPostgres) NewSelectParams(params entity.PaperParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLikeFields([]string{"UPPER(paper_article)", "UPPER(paper_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("type_id", params.Filter.TypeIDs),
				w.FilterAnyOf("color_id", params.Filter.ColorIDs),
				w.FilterAnyOf("facture_id", params.Filter.FactureIDs),
				w.FilterRangeFloat64("paper_width", mrtype.RangeFloat64(params.Filter.Width), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("paper_height", mrtype.RangeFloat64(params.Filter.Height), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("paper_density", mrtype.RangeFloat64(params.Filter.Density), 0, mrlib.EqualityThresholdE9),
				w.FilterAnyOf("paper_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("paper_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *PaperPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Paper, error) {
	whereStr, whereArgs := params.Where.ToSQL()

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
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
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

	rows := make([]entity.Paper, 0)

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

// FetchTotal - comment method.
func (re *PaperPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
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
func (re *PaperPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Paper, error) {
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
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
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
func (re *PaperPostgres) FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			paper_id
		FROM
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
		WHERE
			paper_article = $1 AND deleted_at IS NULL
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
func (re *PaperPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			paper_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
		WHERE
			paper_id = $1 AND deleted_at IS NULL
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
func (re *PaperPostgres) Insert(ctx context.Context, row entity.Paper) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNamePapers + `
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

	err := re.client.Conn(ctx).QueryRow(
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
func (re *PaperPostgres) Update(ctx context.Context, row entity.Paper) (int32, error) {
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
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			paper_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *PaperPostgres) UpdateStatus(ctx context.Context, row entity.Paper) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			paper_status = $3
		WHERE
			paper_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *PaperPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNamePapers + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			paper_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

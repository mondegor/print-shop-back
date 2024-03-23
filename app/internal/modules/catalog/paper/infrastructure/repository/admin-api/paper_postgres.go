package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog/paper"
	entity "print-shop-back/internal/modules/catalog/paper/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewPaperPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *PaperPostgres {
	return &PaperPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *PaperPostgres) NewSelectParams(params entity.PaperParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("paper_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(paper_article)", "UPPER(paper_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("color_id", params.Filter.ColorIDs),
				w.FilterAnyOf("facture_id", params.Filter.FactureIDs),
				w.FilterRangeInt64("paper_length", params.Filter.Length, 0),
				w.FilterRangeInt64("paper_width", params.Filter.Width, 0),
				w.FilterRangeInt64("paper_density", params.Filter.Density, 0),
				w.FilterAnyOf("paper_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("paper_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *PaperPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Paper, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			paper_id,
			tag_version,
			paper_article as article,
			paper_caption as caption,
			color_id,
			facture_id,
			paper_length as length,
			paper_width as width,
			paper_density as density,
			paper_thickness,
			paper_sides,
			paper_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBSchema + `.papers
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

	rows := make([]entity.Paper, 0)

	for cursor.Next() {
		var row entity.Paper

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Article,
			&row.Caption,
			&row.ColorID,
			&row.FactureID,
			&row.Length,
			&row.Width,
			&row.Density,
			&row.Thickness,
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

func (re *PaperPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.papers
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

func (re *PaperPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Paper, error) {
	sql := `
		SELECT
			tag_version,
			paper_article,
			paper_caption,
			color_id,
			facture_id,
			paper_length,
			paper_width,
			paper_density,
			paper_thickness,
			paper_sides,
			paper_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.papers
		WHERE
			paper_id = $1 AND paper_status <> $2
		LIMIT 1;`

	row := entity.Paper{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.Article,
		&row.Caption,
		&row.ColorID,
		&row.FactureID,
		&row.Length,
		&row.Width,
		&row.Density,
		&row.Thickness,
		&row.Sides,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (re *PaperPostgres) FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			paper_id
		FROM
			` + module.DBSchema + `.papers
		WHERE
			paper_article = $1
		LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			paper_status
		FROM
			` + module.DBSchema + `.papers
		WHERE
			paper_id = $1 AND paper_status <> $2
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}

func (re *PaperPostgres) Insert(ctx context.Context, row entity.Paper) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.papers
			(
				paper_article,
				paper_caption,
				color_id,
				facture_id,
				paper_length,
				paper_width,
				paper_density,
				paper_thickness,
				paper_sides,
				paper_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
			paper_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.ColorID,
		row.FactureID,
		row.Length,
		row.Width,
		row.Density,
		row.Thickness,
		row.Sides,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *PaperPostgres) Update(ctx context.Context, row entity.Paper) (int32, error) {
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
			` + module.DBSchema + `.papers
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			paper_id = $1 AND tag_version = $2 AND paper_status <> $3
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

func (re *PaperPostgres) UpdateStatus(ctx context.Context, row entity.Paper) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.papers
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			paper_status = $4
		WHERE
			paper_id = $1 AND tag_version = $2 AND paper_status <> $3
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

func (re *PaperPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.papers
		SET
			tag_version = tag_version + 1,
			paper_article = NULL,
			updated_at = NOW(),
			paper_status = $2
		WHERE
			paper_id = $1 AND paper_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}

package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog/box"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	BoxPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewBoxPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *BoxPostgres {
	return &BoxPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *BoxPostgres) NewSelectParams(params entity.BoxParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("box_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(box_article)", "UPPER(box_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterRangeInt64("box_length", params.Filter.Length, 0),
				w.FilterRangeInt64("box_width", params.Filter.Width, 0),
				w.FilterRangeInt64("box_depth", params.Filter.Depth, 0),
				w.FilterAnyOf("box_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("box_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *BoxPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Box, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			box_id,
			tag_version,
			created_at as createdAt,
			updated_at as updatedAt,
			box_article as article,
			box_caption as caption,
			box_length as length,
			box_width as width,
			box_depth as depth,
			box_status
		FROM
			` + module.DBSchema + `.boxes
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

	rows := make([]entity.Box, 0)

	for cursor.Next() {
		var row entity.Box

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Article,
			&row.Caption,
			&row.Length,
			&row.Width,
			&row.Depth,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *BoxPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.boxes
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

func (re *BoxPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Box, error) {
	sql := `
		SELECT
			tag_version,
			created_at,
			updated_at,
			box_article,
			box_caption,
			box_length,
			box_width,
			box_depth,
			box_status
		FROM
			` + module.DBSchema + `.boxes
		WHERE
			box_id = $1 AND box_status <> $2
		LIMIT 1;`

	row := entity.Box{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.Article,
		&row.Caption,
		&row.Length,
		&row.Width,
		&row.Depth,
		&row.Status,
	)

	return row, err
}

func (re *BoxPostgres) FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			box_id
		FROM
			` + module.DBSchema + `.boxes
		WHERE
			box_article = $1
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

func (re *BoxPostgres) FetchStatus(ctx context.Context, row entity.Box) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			box_status
		FROM
			` + module.DBSchema + `.boxes
		WHERE
			box_id = $1 AND box_status <> $2
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *BoxPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		SELECT
			box_id
		FROM
			` + module.DBSchema + `.boxes
		WHERE
			box_id = $1 AND box_status <> $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}

func (re *BoxPostgres) Insert(ctx context.Context, row entity.Box) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.boxes
			(
				box_article,
				box_caption,
				box_length,
				box_width,
				box_depth,
				box_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			box_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.Length,
		row.Width,
		row.Depth,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *BoxPostgres) Update(ctx context.Context, row entity.Box) (int32, error) {
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
			` + module.DBSchema + `.boxes
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			box_id = $1 AND tag_version = $2 AND box_status <> $3
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

func (re *BoxPostgres) UpdateStatus(ctx context.Context, row entity.Box) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.boxes
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			box_status = $4
		WHERE
			box_id = $1 AND tag_version = $2 AND box_status <> $3
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

func (re *BoxPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.boxes
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			box_article = NULL,
			box_status = $2
		WHERE
			box_id = $1 AND box_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}

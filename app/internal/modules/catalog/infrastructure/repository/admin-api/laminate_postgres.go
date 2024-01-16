package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminatePostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewLaminatePostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *LaminatePostgres {
	return &LaminatePostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *LaminatePostgres) NewFetchParams(params entity.LaminateParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("laminate_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(laminate_article)", "UPPER(laminate_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("type_id", params.Filter.TypeIDs),
				w.FilterRangeInt64("laminate_length", params.Filter.Length, 0),
				w.FilterRangeInt64("laminate_weight", params.Filter.Weight, 0),
				w.FilterAnyOf("laminate_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("laminate_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *LaminatePostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Laminate, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			laminate_id,
			tag_version,
			datetime_created as createdAt,
			datetime_updated as updatedAt,
			laminate_article as article,
			laminate_caption as caption,
			type_id,
			laminate_length as length,
			laminate_weight as weight,
			laminate_thickness,
			laminate_status
		FROM
			` + module.UnitLaminateDBSchema + `.laminates
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

	rows := make([]entity.Laminate, 0)

	for cursor.Next() {
		var row entity.Laminate

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Article,
			&row.Caption,
			&row.TypeID,
			&row.Length,
			&row.Weight,
			&row.Thickness,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *LaminatePostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.UnitLaminateDBSchema + `.laminates
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

func (re *LaminatePostgres) LoadOne(ctx context.Context, row *entity.Laminate) error {
	sql := `
		SELECT
			tag_version,
			datetime_created,
			datetime_updated,
			laminate_article,
			laminate_caption,
			type_id,
			laminate_length,
			laminate_weight,
			laminate_thickness,
			laminate_status
		FROM
			` + module.UnitLaminateDBSchema + `.laminates
		WHERE
			laminate_id = $1 AND laminate_status <> $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.Article,
		&row.Caption,
		&row.TypeID,
		&row.Length,
		&row.Weight,
		&row.Thickness,
		&row.Status,
	)
}

func (re *LaminatePostgres) FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			laminate_id
		FROM
			` + module.UnitLaminateDBSchema + `.laminates
		WHERE
			laminate_article = $1
		LIMIT 1;`

	var id mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&id,
	)

	return id, err
}

func (re *LaminatePostgres) FetchStatus(ctx context.Context, row *entity.Laminate) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			laminate_status
		FROM
			` + module.UnitLaminateDBSchema + `.laminates
		WHERE
			laminate_id = $1 AND laminate_status <> $2
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
func (re *LaminatePostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		SELECT
			1
		FROM
			` + module.UnitBoxDBSchema + `.laminates
		WHERE
			laminate_id = $1 AND laminate_status <> $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	).Scan(
		&id,
	)
}

func (re *LaminatePostgres) Insert(ctx context.Context, row *entity.Laminate) error {
	sql := `
		INSERT INTO ` + module.UnitLaminateDBSchema + `.laminates
			(
				laminate_article,
				laminate_caption,
				type_id,
				laminate_length,
				laminate_weight,
				laminate_thickness,
				laminate_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			laminate_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Article,
		row.Caption,
		row.TypeID,
		row.Length,
		row.Weight,
		row.Thickness,
		row.Status,
	).Scan(
		&row.ID,
	)

	return err
}

func (re *LaminatePostgres) Update(ctx context.Context, row *entity.Laminate) (int32, error) {
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
			` + module.UnitLaminateDBSchema + `.laminates
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
			` + setStr + `
		WHERE
			laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3
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

func (re *LaminatePostgres) UpdateStatus(ctx context.Context, row *entity.Laminate) (int32, error) {
	sql := `
		UPDATE
			` + module.UnitLaminateDBSchema + `.laminates
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
			laminate_status = $4
		WHERE
			laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3
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

func (re *LaminatePostgres) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.UnitLaminateDBSchema + `.laminates
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
			laminate_article = NULL,
			laminate_status = $2
		WHERE
			laminate_id = $1 AND laminate_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	)
}

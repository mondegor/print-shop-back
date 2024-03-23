package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-facture"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/admin-api"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-facture/infrastructure/repository/shared"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFacturePostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewPaperFacturePostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *PaperFacturePostgres) NewSelectParams(params entity.PaperFactureParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("facture_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(facture_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("facture_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("facture_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *PaperFacturePostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PaperFacture, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            facture_id,
            tag_version,
            facture_caption as caption,
            facture_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.paper_factures
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

	rows := make([]entity.PaperFacture, 0)

	for cursor.Next() {
		var row entity.PaperFacture

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Caption,
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

func (re *PaperFacturePostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.paper_factures
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

func (re *PaperFacturePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperFacture, error) {
	sql := `
        SELECT
            tag_version,
            facture_caption,
            facture_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.paper_factures
        WHERE
            facture_id = $1 AND facture_status <> $2
        LIMIT 1;`

	row := entity.PaperFacture{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.PaperFactureFetchStatusPostgres(ctx, re.client, rowID)
}

func (re *PaperFacturePostgres) Insert(ctx context.Context, row entity.PaperFacture) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.paper_factures
            (
                facture_caption,
                facture_status
            )
        VALUES
            ($1, $2)
        RETURNING
            facture_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *PaperFacturePostgres) Update(ctx context.Context, row entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_caption = $4
        WHERE
            facture_id = $1 AND tag_version = $2 AND facture_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *PaperFacturePostgres) UpdateStatus(ctx context.Context, row entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_status = $4
        WHERE
            facture_id = $1 AND tag_version = $2 AND facture_status <> $3
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

func (re *PaperFacturePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_status = $2
        WHERE
            facture_id = $1 AND facture_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}

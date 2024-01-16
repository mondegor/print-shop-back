package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"
	repository_shared "print-shop-back/internal/modules/dictionaries/infrastructure/repository/shared"
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

func (re *PaperFacturePostgres) NewFetchParams(params entity.PaperFactureParams) mrstorage.SqlSelectParams {
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
            datetime_created as createdAt,
			datetime_updated as updatedAt,
            facture_caption as caption,
            facture_status
        FROM
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
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
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Caption,
			&row.Status,
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
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
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

func (re *PaperFacturePostgres) LoadOne(ctx context.Context, row *entity.PaperFacture) error {
	sql := `
        SELECT
            tag_version,
            datetime_created,
			datetime_updated,
            facture_caption,
            facture_status
        FROM
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
        WHERE
            facture_id = $1 AND facture_status <> $2
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
		&row.Caption,
		&row.Status,
	)
}

func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, row *entity.PaperFacture) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            facture_status
        FROM
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
        WHERE
            facture_id = $1 AND facture_status <> $2
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
func (re *PaperFacturePostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.PaperFactureIsExistsPostgres(ctx, re.client, id)
}

func (re *PaperFacturePostgres) Insert(ctx context.Context, row *entity.PaperFacture) error {
	sql := `
        INSERT INTO ` + module.UnitPaperFactureDBSchema + `.paper_factures
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

	return err
}

func (re *PaperFacturePostgres) Update(ctx context.Context, row *entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *PaperFacturePostgres) UpdateStatus(ctx context.Context, row *entity.PaperFacture) (int32, error) {
	sql := `
        UPDATE
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *PaperFacturePostgres) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.UnitPaperFactureDBSchema + `.paper_factures
        SET
            tag_version = tag_version + 1,
			datetime_updated = NOW(),
            facture_status = $2
        WHERE
            facture_id = $1 AND facture_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	)
}

package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// PaperFacturePostgres - comment struct.
	PaperFacturePostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewPaperFacturePostgres - создаёт объект PaperFacturePostgres.
func NewPaperFacturePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNamePaperFactures,
			"facture_id",
			module.DBFieldTagVersion,
			"facture_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNamePaperFactures,
			"facture_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNamePaperFactures,
		),
	}
}

// FetchWithTotal - comment method.
func (re *PaperFacturePostgres) FetchWithTotal(
	ctx context.Context,
	params entity.PaperFactureParams,
) (rows []entity.PaperFacture, countRows uint64, err error) {
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
func (re *PaperFacturePostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.PaperFacture, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
        SELECT
            facture_id,
            tag_version,
            facture_caption as caption,
            facture_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBTableNamePaperFactures + `
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

	rows := make([]entity.PaperFacture, 0, maxRows)

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

// FetchOne - comment method.
func (re *PaperFacturePostgres) FetchOne(ctx context.Context, rowID uint64) (entity.PaperFacture, error) {
	sql := `
        SELECT
            tag_version,
            facture_caption,
            facture_status,
            created_at,
			updated_at
        FROM
            ` + module.DBTableNamePaperFactures + `
        WHERE
            facture_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.PaperFacture{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *PaperFacturePostgres) Insert(ctx context.Context, row entity.PaperFacture) (rowID uint64, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNamePaperFactures + `
            (
                facture_caption,
                facture_status
            )
        VALUES
            ($1, $2)
        RETURNING
            facture_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *PaperFacturePostgres) Update(ctx context.Context, row entity.PaperFacture) (tagVersion uint32, err error) {
	sql := `
        UPDATE
            ` + module.DBTableNamePaperFactures + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            facture_caption = $3
        WHERE
            facture_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *PaperFacturePostgres) UpdateStatus(ctx context.Context, row entity.PaperFacture) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *PaperFacturePostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *PaperFacturePostgres) fetchCondition(filter entity.PaperFactureListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(facture_caption)", strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("facture_status", filter.Statuses),
			)
		},
	)
}

func (re *PaperFacturePostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("facture_id", mrenum.SortDirectionASC),
			)
		},
	)
}

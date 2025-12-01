package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/measure"
)

type (
	// PaperPostgres - comment struct.
	PaperPostgres struct {
		client         mrstorage.DBConnManager
		repoTypeIDs    db.ColumnFetcher[itemstatus.Enum, uint64]
		repoColorIDs   db.ColumnFetcher[itemstatus.Enum, uint64]
		repoDensities  db.ColumnFetcher[itemstatus.Enum, measure.KilogramPerMeter2]
		repoFactureIDs db.ColumnFetcher[itemstatus.Enum, uint64]
	}
)

// NewPaperPostgres - создаёт объект PaperPostgres.
func NewPaperPostgres(client mrstorage.DBConnManager) *PaperPostgres {
	return &PaperPostgres{
		client: client,
		repoTypeIDs: db.NewColumnFetcher[itemstatus.Enum, uint64](
			client,
			module.DBTableNamePapers,
			"paper_status",
			"type_id",
			module.DBFieldDeletedAt,
		),
		repoColorIDs: db.NewColumnFetcher[itemstatus.Enum, uint64](
			client,
			module.DBTableNamePapers,
			"paper_status",
			"color_id",
			module.DBFieldDeletedAt,
		),
		repoDensities: db.NewColumnFetcher[itemstatus.Enum, measure.KilogramPerMeter2](
			client,
			module.DBTableNamePapers,
			"paper_status",
			"paper_density",
			module.DBFieldDeletedAt,
		),
		repoFactureIDs: db.NewColumnFetcher[itemstatus.Enum, uint64](
			client,
			module.DBTableNamePapers,
			"paper_status",
			"facture_id",
			module.DBFieldDeletedAt,
		),
	}
}

// Fetch - comment method.
func (re *PaperPostgres) Fetch(ctx context.Context, _ entity.PaperParams) ([]entity.Paper, error) {
	sql := `
        SELECT
            paper_id,
			paper_article,
			paper_caption,
			type_id,
			color_id,
			facture_id,
			paper_width,
			paper_height,
			paper_thickness,
			paper_density,
			paper_sides
        FROM
            ` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        ORDER BY
            paper_caption ASC, paper_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		itemstatus.Enabled,
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
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// FetchTypeIDs - comment method.
func (re *PaperPostgres) FetchTypeIDs(ctx context.Context) ([]uint64, error) {
	return re.repoTypeIDs.Fetch(ctx, itemstatus.Enabled)
}

// FetchColorIDs - comment method.
func (re *PaperPostgres) FetchColorIDs(ctx context.Context) ([]uint64, error) {
	return re.repoColorIDs.Fetch(ctx, itemstatus.Enabled)
}

// FetchDensities - comment method.
func (re *PaperPostgres) FetchDensities(ctx context.Context) ([]measure.KilogramPerMeter2, error) {
	return re.repoDensities.Fetch(ctx, itemstatus.Enabled)
}

// FetchFactureIDs - comment method.
func (re *PaperPostgres) FetchFactureIDs(ctx context.Context) ([]uint64, error) {
	return re.repoFactureIDs.Fetch(ctx, itemstatus.Enabled)
}

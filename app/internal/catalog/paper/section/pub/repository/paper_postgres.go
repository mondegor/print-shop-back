package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// PaperPostgres - comment struct.
	PaperPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewPaperPostgres - создаёт объект PaperPostgres.
func NewPaperPostgres(client mrstorage.DBConnManager) *PaperPostgres {
	return &PaperPostgres{
		client: client,
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
            ` + module.DBSchema + `.` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        ORDER BY
            paper_caption ASC, paper_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
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
func (re *PaperPostgres) FetchTypeIDs(ctx context.Context) ([]mrtype.KeyInt32, error) {
	sql := `
        SELECT
			type_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        GROUP BY
            type_id
		ORDER BY type_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]mrtype.KeyInt32, 0)

	for cursor.Next() {
		var typeID mrtype.KeyInt32

		err = cursor.Scan(
			&typeID,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, typeID)
	}

	return rows, cursor.Err()
}

// FetchColorIDs - comment method.
func (re *PaperPostgres) FetchColorIDs(ctx context.Context) ([]mrtype.KeyInt32, error) {
	sql := `
        SELECT
			color_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        GROUP BY
            color_id
		ORDER BY color_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]mrtype.KeyInt32, 0)

	for cursor.Next() {
		var colorID mrtype.KeyInt32

		err = cursor.Scan(
			&colorID,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, colorID)
	}

	return rows, cursor.Err()
}

// FetchDensities - comment method.
func (re *PaperPostgres) FetchDensities(ctx context.Context) ([]measure.KilogramPerMeter2, error) {
	sql := `
        SELECT
			paper_density
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        GROUP BY
            paper_density
		ORDER BY paper_density ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]measure.KilogramPerMeter2, 0)

	for cursor.Next() {
		var typeID measure.KilogramPerMeter2

		err = cursor.Scan(
			&typeID,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, typeID)
	}

	return rows, cursor.Err()
}

// FetchFactureIDs - comment method.
func (re *PaperPostgres) FetchFactureIDs(ctx context.Context) ([]mrtype.KeyInt32, error) {
	sql := `
        SELECT
			facture_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePapers + `
        WHERE
            paper_status = $1 AND deleted_at IS NULL
        GROUP BY
            facture_id
		ORDER BY facture_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]mrtype.KeyInt32, 0)

	for cursor.Next() {
		var factureID mrtype.KeyInt32

		err = cursor.Scan(
			&factureID,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, factureID)
	}

	return rows, cursor.Err()
}

package repository

import (
	"context"
	"strings"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// CompanyPagePostgres - comment struct.
	CompanyPagePostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewCompanyPagePostgres - создаёт объект CompanyPagePostgres.
func NewCompanyPagePostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *CompanyPagePostgres) NewSelectParams(params entity.CompanyPageParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(page_title)", "UPPER(site_url)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("page_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("account_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *CompanyPagePostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.CompanyPage, error) {
	whereStr, whereArgs := params.Where.WithPrefix(" WHERE ").ToSQL()

	sql := `
        SELECT
            account_id,
            rewrite_name as rewriteName,
            page_title as pageTitle,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url as siteUrl,
            page_status,
			created_at as createdAt,
            updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
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

	rows := make([]entity.CompanyPage, 0)

	for cursor.Next() {
		var row entity.CompanyPage

		err = cursor.Scan(
			&row.AccountID,
			&row.RewriteName,
			&row.PageTitle,
			&row.LogoURL,
			&row.SiteURL,
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
func (re *CompanyPagePostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.WithPrefix(" WHERE ").ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
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

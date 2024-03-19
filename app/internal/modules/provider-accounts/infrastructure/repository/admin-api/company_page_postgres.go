package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	CompanyPagePostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewCompanyPagePostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *CompanyPagePostgres) NewSelectParams(params entity.CompanyPageParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(page_head)", "UPPER(site_url)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("page_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("account_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *CompanyPagePostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CompanyPage, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            account_id,
            updated_at as updatedAt,
            rewrite_name as rewriteName,
            page_head as pageHead,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url as siteUrl,
            page_status
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
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

	rows := make([]entity.CompanyPage, 0)

	for cursor.Next() {
		var row entity.CompanyPage

		err = cursor.Scan(
			&row.AccountID,
			&row.UpdatedAt,
			&row.RewriteName,
			&row.PageHead,
			&row.LogoURL,
			&row.SiteURL,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *CompanyPagePostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
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

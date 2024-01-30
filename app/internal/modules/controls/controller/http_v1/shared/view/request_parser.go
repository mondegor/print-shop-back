package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	RequestParser interface {
		mrserver.RequestParserInt64
		mrserver.RequestParserItemStatus
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserSortPage
		mrserver.RequestParserString
		mrserver.RequestParserValidate

		FormDataID(r *http.Request) (mrtype.KeyInt32, error)
		FilterElementDetailingList(r *http.Request, key string) []entity_shared.ElementDetailing
	}

	Parser struct {
		*mrparser.Int64
		*mrparser.ItemStatus
		*mrparser.KeyInt32
		*mrparser.SortPage
		*mrparser.String
		*mrparser.Validator
	}
)

func NewParser(
	p1 *mrparser.Int64,
	p2 *mrparser.ItemStatus,
	p3 *mrparser.KeyInt32,
	p4 *mrparser.SortPage,
	p5 *mrparser.String,
	p6 *mrparser.Validator,
) *Parser {
	return &Parser{
		Int64:      p1,
		ItemStatus: p2,
		KeyInt32:   p3,
		SortPage:   p4,
		String:     p5,
		Validator:  p6,
	}
}

func (p *Parser) FormDataID(r *http.Request) (mrtype.KeyInt32, error) {
	value, err := mrreq.ParseInt64(r, "formId", true)

	if err != nil {
		return 0, usecase_shared.FactoryErrFormDataNotFound.Wrap(err, p.RawParamString(r, "formId"))
	}

	return mrtype.KeyInt32(value), nil
}

func (p *Parser) FilterElementDetailingList(r *http.Request, key string) []entity_shared.ElementDetailing {
	items, err := parseFilterDetailingList(
		r,
		key,
		[]entity_shared.ElementDetailing{
			entity_shared.ElementDetailingNormal,
		},
	)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()
	}

	return items
}

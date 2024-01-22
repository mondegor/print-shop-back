package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	RequestParser interface {
		mrserver.RequestParser
		mrserver.RequestParserItemStatus
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserSortPage
		mrserver.RequestParserValidate

		FormDataID(r *http.Request) (mrtype.KeyInt32, error)
		FilterElementDetailingList(r *http.Request, key string) []entity_shared.ElementDetailing
	}

	Parser struct {
		*mrparser.Base
		*mrparser.ItemStatus
		*mrparser.KeyInt32
		*mrparser.SortPage
		*mrparser.Validator
	}
)

func NewParser(
	p1 *mrparser.Base,
	p2 *mrparser.ItemStatus,
	p3 *mrparser.KeyInt32,
	p4 *mrparser.SortPage,
	p5 *mrparser.Validator,
) *Parser {
	return &Parser{
		Base:       p1,
		ItemStatus: p2,
		KeyInt32:   p3,
		SortPage:   p4,
		Validator:  p5,
	}
}

func (p *Parser) FormDataID(r *http.Request) (mrtype.KeyInt32, error) {
	value, err := mrreq.ParseInt64(r, "formId", true)

	if err != nil {
		return 0, usecase_shared.FactoryErrFormDataNotFound.Wrap(err, p.RawQueryParam(r, "formId"))
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
		mrctx.Logger(r.Context()).Warn(err)
	}

	return items
}

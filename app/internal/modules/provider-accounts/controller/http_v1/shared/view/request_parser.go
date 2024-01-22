package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	RequestParser interface {
		mrserver.RequestParser
		mrserver.RequestParserItemStatus
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserSortPage
		mrserver.RequestParserValidate

		FilterPublicStatusList(r *http.Request, key string) []entity_shared.PublicStatus
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

func (p *Parser) FilterPublicStatusList(r *http.Request, key string) []entity_shared.PublicStatus {
	items, err := parseFilterPublicStatusList(
		r,
		key,
		[]entity_shared.PublicStatus{
			entity_shared.PublicStatusPublished,
			entity_shared.PublicStatusPublishedShared,
		},
	)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
	}

	return items
}

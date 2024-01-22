package view_shared

import (
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

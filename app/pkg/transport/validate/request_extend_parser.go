package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// RequestExtendParser - comment interface.
	RequestExtendParser interface {
		RequestParser
		request.ParserItemStatus
		request.ParserListPager
		request.ParserListSorter
	}

	// ExtendParser - comment struct.
	ExtendParser struct {
		*Parser
		*parser.ItemStatus
		*parser.ListPager
		*parser.ListSorter
	}
)

// NewExtendParser - создаёт объект Parser.
func NewExtendParser(
	p1 *Parser,
	p2 *parser.ItemStatus,
	p3 *parser.ListPager,
	p4 *parser.ListSorter,
) *ExtendParser {
	return &ExtendParser{
		Parser:     p1,
		ItemStatus: p2,
		ListPager:  p3,
		ListSorter: p4,
	}
}

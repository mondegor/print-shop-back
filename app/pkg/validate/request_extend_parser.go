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
		request.ParserListSorter
		request.ParserListPager
	}

	// ExtendParser - comment struct.
	ExtendParser struct {
		*Parser
		*parser.ItemStatus
		*parser.ListSorter
		*parser.ListPager
	}
)

// NewExtendParser - создаёт объект Parser.
func NewExtendParser(
	p1 *Parser,
	p2 *parser.ItemStatus,
	p3 *parser.ListSorter,
	p4 *parser.ListPager,
) *ExtendParser {
	return &ExtendParser{
		Parser:     p1,
		ItemStatus: p2,
		ListSorter: p3,
		ListPager:  p4,
	}
}

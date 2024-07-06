package validate

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	// RequestExtendParser - comment interface.
	RequestExtendParser interface {
		RequestParser
		mrserver.RequestParserItemStatus
		mrserver.RequestParserListSorter
		mrserver.RequestParserListPager
	}

	// ExtendParser - comment struct.
	ExtendParser struct {
		*Parser
		*mrparser.ItemStatus
		*mrparser.ListSorter
		*mrparser.ListPager
	}
)

// NewExtendParser - создаёт объект Parser.
func NewExtendParser(
	p1 *Parser,
	p2 *mrparser.ItemStatus,
	p3 *mrparser.ListSorter,
	p4 *mrparser.ListPager,
) *ExtendParser {
	return &ExtendParser{
		Parser:     p1,
		ItemStatus: p2,
		ListSorter: p3,
		ListPager:  p4,
	}
}

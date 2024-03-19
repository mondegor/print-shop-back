package view_shared

import (
	"print-shop-back/pkg/modules/controls/view"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	RequestParser interface {
		mrserver.RequestParserInt64
		mrserver.RequestParserItemStatus
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserListSorter
		mrserver.RequestParserListPager
		mrserver.RequestParserString
		mrserver.RequestParserValidate
		mrserver.RequestParserFile
		// mrserver.RequestParserImage
		view.RequestDetailingParser
	}

	Parser struct {
		*mrparser.Int64
		*mrparser.ItemStatus
		*mrparser.KeyInt32
		*mrparser.ListSorter
		*mrparser.ListPager
		*mrparser.String
		*mrparser.Validator
		*mrparser.File
		*view.DetailingParser
	}
)

func NewParser(
	p1 *mrparser.Int64,
	p2 *mrparser.ItemStatus,
	p3 *mrparser.KeyInt32,
	p4 *mrparser.ListSorter,
	p5 *mrparser.ListPager,
	p6 *mrparser.String,
	p7 *mrparser.Validator,
	p8 *mrparser.File,
	p9 *view.DetailingParser,
) *Parser {
	return &Parser{
		Int64:           p1,
		ItemStatus:      p2,
		KeyInt32:        p3,
		ListSorter:      p4,
		ListPager:       p5,
		String:          p6,
		Validator:       p7,
		File:            p8,
		DetailingParser: p9,
	}
}

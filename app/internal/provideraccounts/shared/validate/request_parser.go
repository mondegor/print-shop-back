package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/transport/validate"
)

type (
	// RequestProviderAccountsParser - comment interface.
	RequestProviderAccountsParser interface {
		pkgvalidate.RequestExtendParser
		request.ParserUser
		request.ParserImage
		validate.RequestPublicStatusParser
	}

	// Parser - comment struct.
	Parser struct {
		*pkgvalidate.ExtendParser
		*parser.User
		*parser.Image
		*validate.PublicStatusParser
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *pkgvalidate.ExtendParser,
	p2 *parser.User,
	p3 *parser.Image,
	p4 *validate.PublicStatusParser,
) *Parser {
	return &Parser{
		ExtendParser:       p1,
		User:               p2,
		Image:              p3,
		PublicStatusParser: p4,
	}
}

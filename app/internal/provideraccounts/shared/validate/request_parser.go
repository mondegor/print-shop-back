package validate

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/validate"
)

type (
	// RequestProviderAccountsParser - comment interface.
	RequestProviderAccountsParser interface {
		pkgvalidate.RequestExtendParser
		mrserver.RequestParserUser
		mrserver.RequestParserImage
		validate.RequestPublicStatusParser
	}

	// Parser - comment struct.
	Parser struct {
		*pkgvalidate.ExtendParser
		*mrparser.User
		*mrparser.Image
		*validate.PublicStatusParser
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *pkgvalidate.ExtendParser,
	p2 *mrparser.User,
	p3 *mrparser.Image,
	p4 *validate.PublicStatusParser,
) *Parser {
	return &Parser{
		ExtendParser:       p1,
		User:               p2,
		Image:              p3,
		PublicStatusParser: p4,
	}
}

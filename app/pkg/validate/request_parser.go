package validate

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	// RequestParser - comment interface.
	RequestParser interface {
		mrserver.RequestParserInt64
		mrserver.RequestParserUint64
		mrserver.RequestParserString
		mrserver.RequestParserUUID
		mrserver.RequestParserValidate
		mrserver.RequestParserClientIP
		mrserver.RequestParserUser
		mrserver.RequestParserLocale
	}

	// Parser - comment struct.
	Parser struct {
		*mrparser.Int64
		*mrparser.Uint64
		*mrparser.String
		*mrparser.UUID
		*mrparser.Validator
		*mrparser.ClientIP
		*mrparser.User
		*mrparser.Locale
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *mrparser.Int64,
	p2 *mrparser.Uint64,
	p3 *mrparser.String,
	p4 *mrparser.UUID,
	p5 *mrparser.Validator,
	p6 *mrparser.ClientIP,
	p7 *mrparser.User,
	p8 *mrparser.Locale,
) *Parser {
	return &Parser{
		Int64:     p1,
		Uint64:    p2,
		String:    p3,
		UUID:      p4,
		Validator: p5,
		ClientIP:  p6,
		User:      p7,
		Locale:    p8,
	}
}

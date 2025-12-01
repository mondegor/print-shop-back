package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// RequestParser - comment interface.
	RequestParser interface {
		request.ParserInt64
		request.ParserUint64
		request.ParserString
		request.ParserUUID
		request.ParserValidate
		request.ParserClientIP
		request.ParserUser
		request.ParserLocale
	}

	// Parser - comment struct.
	Parser struct {
		*parser.Int64
		*parser.Uint64
		*parser.String
		*parser.UUID
		*parser.Validator
		*parser.ClientIP
		*parser.User
		*parser.Locale
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *parser.Int64,
	p2 *parser.Uint64,
	p3 *parser.String,
	p4 *parser.UUID,
	p5 *parser.Validator,
	p6 *parser.ClientIP,
	p7 *parser.User,
	p8 *parser.Locale,
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

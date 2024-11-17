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
	}

	// Parser - comment struct.
	Parser struct {
		*mrparser.Int64
		*mrparser.Uint64
		*mrparser.String
		*mrparser.UUID
		*mrparser.Validator
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *mrparser.Int64,
	p2 *mrparser.Uint64,
	p3 *mrparser.String,
	p4 *mrparser.UUID,
	p5 *mrparser.Validator,
) *Parser {
	return &Parser{
		Int64:     p1,
		Uint64:    p2,
		String:    p3,
		UUID:      p4,
		Validator: p5,
	}
}

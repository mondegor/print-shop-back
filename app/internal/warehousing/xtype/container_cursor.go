package xtype

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// ContainerCursor - параметры для выборки части списка элементов.
	ContainerCursor struct {
		Code   string
		Marker uint16
		Limit  int
	}
)

// NewContainerCursor - создаёт объект ContainerCursor.
func NewContainerCursor(params mrtype.CursorParams) ContainerCursor {
	if params.Limit == 0 {
		params.Limit = 1
	}

	code, marker, found := strings.Cut(params.Value, "|")
	if !found {
		return ContainerCursor{
			Limit: params.Limit,
		}
	}

	parsedMarker, err := strconv.ParseUint(marker, 10, 16)
	if err != nil {
		return ContainerCursor{
			Code:  code,
			Limit: params.Limit,
		}
	}

	return ContainerCursor{
		Code:   code,
		Marker: uint16(parsedMarker),
		Limit:  params.Limit,
	}
}

package api

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/controls/enum"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ElementTemplateHeaderName = "Controls.API.ElementTemplateHeader" // ElementTemplateHeaderName - название API
)

type (
	// ElementTemplateDTO - comment struct.
	ElementTemplateDTO struct {
		ID         mrtype.KeyInt32
		TagVersion int32
		ParamName  string
		Caption    string
		Detailing  enum.ElementDetailing
	}

	// ElementTemplateHeader - comment interface.
	ElementTemplateHeader interface {
		// GetItemHeader - error: ErrElementTemplateRequired |
		//                        ErrElementTemplateNotFound |
		//                        ErrElementTemplateIsDisabled |
		//                        Failed
		GetItemHeader(ctx context.Context, itemID mrtype.KeyInt32) (ElementTemplateDTO, error)
	}
)

var (
	// ErrElementTemplateRequired - element template ID is required.
	ErrElementTemplateRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsElementTemplateRequired", mrerr.ErrorKindUser, "element template ID is required")

	// ErrElementTemplateNotFound - element template with ID not found.
	ErrElementTemplateNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsElementTemplateNotFound", mrerr.ErrorKindUser, "element template with ID={{ .id }} not found")

	// ErrElementTemplateIsDisabled - element template with ID is disabled.
	ErrElementTemplateIsDisabled = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsElementTemplateIsDisabled", mrerr.ErrorKindUser, "element template with ID={{ .id }} is disabled")
)

// ElementTemplateErrors - comment func.
func ElementTemplateErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrElementTemplateRequired,
		ErrElementTemplateNotFound,
		ErrElementTemplateIsDisabled,
	}
}

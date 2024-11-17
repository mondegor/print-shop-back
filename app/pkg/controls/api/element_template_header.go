package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

const (
	ElementTemplateHeaderName = "Controls.API.ElementTemplateHeader" // ElementTemplateHeaderName - название API
)

type (
	// ElementTemplateDTO - comment struct.
	ElementTemplateDTO struct {
		ID         uint64
		TagVersion uint32
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
		GetItemHeader(ctx context.Context, itemID uint64) (ElementTemplateDTO, error)
	}
)

var (
	// ErrElementTemplateRequired - element template ID is required.
	ErrElementTemplateRequired = mrerr.NewProto(
		"controls.errElementTemplateRequired", mrerr.ErrorKindUser, "element template ID is required")

	// ErrElementTemplateNotFound - element template with ID not found.
	ErrElementTemplateNotFound = mrerr.NewProto(
		"controls.errElementTemplateNotFound", mrerr.ErrorKindUser, "element template with ID={{ .id }} not found")

	// ErrElementTemplateIsDisabled - element template with ID is disabled.
	ErrElementTemplateIsDisabled = mrerr.NewProto(
		"controls.errElementTemplateIsDisabled", mrerr.ErrorKindUser, "element template with ID={{ .id }} is disabled")
)

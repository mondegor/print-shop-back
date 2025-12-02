package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/controls/enum/elementdetailing"
)

const (
	// ElementTemplateHeaderName - название API.
	ElementTemplateHeaderName = "Controls.API.ElementTemplateHeader"
)

type (
	// ElementTemplateDTO - comment struct.
	ElementTemplateDTO struct {
		ID         uint64
		TagVersion uint32
		ParamName  string
		Caption    string
		Detailing  elementdetailing.Enum
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
	ErrElementTemplateRequired = mrerr.NewKindUser("ElementTemplateRequired", "element template ID is required")

	// ErrElementTemplateNotFound - element template with ID not found.
	ErrElementTemplateNotFound = mrerr.NewKindUser("ElementTemplateNotFound", "element template with ID={Id} not found")

	// ErrElementTemplateIsDisabled - element template with ID is disabled.
	ErrElementTemplateIsDisabled = mrerr.NewKindUser("ElementTemplateIsDisabled", "element template with ID={Id} is disabled")
)

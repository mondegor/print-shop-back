package controls

import (
	"context"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ElementTemplateAPIName = "Controls.ElementTemplateAPI"
)

type (
	ElementTemplateHead struct {
		ID         mrtype.KeyInt32
		TagVersion int32
		ParamName  string
		Caption    string
		Detailing  enums.ElementDetailing
	}

	ElementTemplateAPI interface {
		// GetItemHead - error: FactoryErrElementTemplateRequired |
		//                      FactoryErrElementTemplateNotFound |
		//                      FactoryErrElementTemplateIsDisabled |
		//                      Failed
		GetItemHead(ctx context.Context, itemID mrtype.KeyInt32) (ElementTemplateHead, error)
	}
)

var (
	FactoryErrElementTemplateRequired = mrerr.NewFactory(
		"errControlsElementTemplateRequired", mrerr.ErrorKindUser, "element template ID is required")

	FactoryErrElementTemplateNotFound = mrerr.NewFactory(
		"errControlsElementTemplateNotFound", mrerr.ErrorKindUser, "element template with ID={{ .id }} not found")

	FactoryErrElementTemplateIsDisabled = mrerr.NewFactory(
		"errControlsElementTemplateIsDisabled", mrerr.ErrorKindUser, "element template with ID={{ .id }} is disabled")
)

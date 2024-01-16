package entity_api

import (
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateHead struct { // DB: ps_controls.element_templates
		ID        mrtype.KeyInt32
		ParamName string
		Caption   string
	}
)

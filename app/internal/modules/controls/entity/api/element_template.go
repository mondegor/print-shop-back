package entity_api

import (
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateHead struct { // DB: printshop_controls.element_templates
		ID        mrtype.KeyInt32
		ParamName string
		Caption   string
	}
)

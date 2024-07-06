package uiform

const (
	ModelNameUIForm      = "uiform.UIForm"      // ModelNameUIForm - название сущности
	ModelNameUIFieldItem = "uiform.UIFieldItem" // ModelNameUIFieldItem - название сущности
)

type (
	// UIForm - comment struct.
	UIForm struct {
		ID      string        `json:"id"`
		Caption string        `json:"caption"`
		Fields  []UIFieldItem `json:"fields"`
	}

	// UIFieldItem - comment struct.
	UIFieldItem struct {
		ID            string        `json:"id"`                      // for all
		Caption       string        `json:"caption,omitempty"`       // for field, enum
		Type          UIDataType    `json:"type,omitempty"`          // for field
		IsRequired    *bool         `json:"required,omitempty"`      // for field
		View          UIItemView    `json:"view,omitempty"`          // for field
		EnabledValues []UIFieldItem `json:"enabledValues,omitempty"` // for field
		Values        []UIFieldItem `json:"values,omitempty"`        // for field
		Measure       string        `json:"measure,omitempty"`       // for field

		IsChecked *bool `json:"checked,omitempty"` // for boolean

		DefaultValue *UIMixedValue `json:"defaultValue,omitempty"` // for number, string
		MinValue     float64       `json:"minValue,omitempty"`     // for number
		MaxValue     float64       `json:"maxValue,omitempty"`     // for number
		MinLength    int32         `json:"minLength,omitempty"`    // for number, string
		MaxLength    int32         `json:"maxLength,omitempty"`    // for number, string
	}
)

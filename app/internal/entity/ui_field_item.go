package entity

type (
    UIForm struct {
        Id     string        `json:"id"`
        Name   string        `json:"name"`
        Fields []UIFieldItem `json:"fields"`
    }

    UIFieldItem struct {
        Id            string        `json:"id"`                      // for all
        Name          string        `json:"name,omitempty"`          // for field, enum
        Type          UIDataType    `json:"type,omitempty"`          // for field
        IsRequired    *bool         `json:"required,omitempty"`      // for field
        View          UIItemView    `json:"view,omitempty"`          // for field
        EnabledValues []UIFieldItem `json:"enabledValues,omitempty"` // for field
        Values        []UIFieldItem `json:"values,omitempty"`        // for field
        Unit          string        `json:"unit,omitempty"`          // for field

        IsChecked *bool `json:"checked,omitempty"` // for boolean

        DefaultValue *UIMixedValue `json:"defaultValue,omitempty"` // for number, string
        MinValue float64           `json:"minValue,omitempty"` // for number
        MaxValue float64           `json:"maxValue,omitempty"` // for number
        MinLength int32            `json:"minLength,omitempty"` // for number, string
        MaxLength int32            `json:"maxLength,omitempty"` // for number, string
    }
)

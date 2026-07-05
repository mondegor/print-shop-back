package entity

const (
	// ModelNameMaterialType - название сущности.
	ModelNameMaterialType = "public-api.Dictionaries.MaterialType"
)

type (
	// MaterialType - comment struct.
	MaterialType struct { // DB: printshop_dictionaries.material_types
		ID      uint64 `json:"id"` // type_id
		Caption string `json:"caption"`
	}

	// MaterialTypeParams - comment struct.
	MaterialTypeParams struct{}
)

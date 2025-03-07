package model

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

type (
	// SheetImpositionRequest - comment struct.
	SheetImpositionRequest struct {
		ItemFormat      string `json:"itemFormat" validate:"required,max=16,tag_2d_size"` // mm x mm
		ItemDistance    string `json:"itemDistance" validate:"max=16,tag_2d_size"`        // mm x mm
		OutFormat       string `json:"outFormat" validate:"required,max=16,tag_2d_size"`  // mm x mm
		DisableRotation bool   `json:"disableRotation"`
		UseMirror       bool   `json:"useMirror"`
	}

	// SheetImpositionResponse - результат вычислений спуска полос.
	SheetImpositionResponse struct {
		ContainerFormat  rect2d.Format     `json:"containerFormat"`
		FragmentDistance measure.Meter     `json:"fragmentDistance"`
		Fragments        []rect2d.Fragment `json:"fragments"`
		TotalElements    uint64            `json:"totalElements"`
		Garbage          measure.Meter2    `json:"garbage"`
		AllowRotation    bool              `json:"allowRotation"`
		UseMirror        bool              `json:"useMirror"`
	}

	// SheetImpositionVariantsResponse - comment struct.
	SheetImpositionVariantsResponse []SheetImpositionResponse
)

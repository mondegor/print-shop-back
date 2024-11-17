package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// CreatePaperFactureRequest - comment struct.
	CreatePaperFactureRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	// StorePaperFactureRequest - comment struct.
	StorePaperFactureRequest struct {
		TagVersion uint32 `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=64"`
	}

	// PaperFactureListResponse - comment struct.
	PaperFactureListResponse struct {
		Items []entity.PaperFacture `json:"items"`
		Total uint64                `json:"total"`
	}
)

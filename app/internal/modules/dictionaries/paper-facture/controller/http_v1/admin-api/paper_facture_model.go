package http_v1

import entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/admin-api"

type (
	CreatePaperFactureRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	StorePaperFactureRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=64"`
	}

	PaperFactureListResponse struct {
		Items []entity.PaperFacture `json:"items"`
		Total int64                 `json:"total"`
	}
)

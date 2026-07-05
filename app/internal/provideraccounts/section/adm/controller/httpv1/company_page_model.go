package httpv1

import (
	"print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPageListResponse - comment struct.
	CompanyPageListResponse struct {
		Items []entity.CompanyPage `json:"items"`
		Total int                  `json:"total"`
	}
)

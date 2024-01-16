package view

import entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"

type (
	CompanyPageListResponse struct {
		Items []entity.CompanyPage `json:"items"`
		Total int64                `json:"total"`
	}
)

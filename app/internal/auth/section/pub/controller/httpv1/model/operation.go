package model

import (
	"github.com/mondegor/go-components/mrauth/enum"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

type (
	// OperationTokenRequest - запрос на подтверждение операции.
	OperationTokenRequest struct {
		Token string `json:"token" validate:"required,min=64,max=128"`
	}

	// ConfirmOperationRequest - запрос на подтверждение операции.
	ConfirmOperationRequest struct {
		Token  string `json:"token" validate:"required,min=64,max=128"`
		Secret string `json:"secret" validate:"required,min=4,max=32"`
	}

	// WaitingConfirmOperationResponse - информация для подтверждения операции.
	WaitingConfirmOperationResponse struct {
		Token             string             `json:"token"`
		ConfirmMethod     enum.ConfirmMethod `json:"confirmMethod"`
		RemainingAttempts uint32             `json:"remainingAttempts"`
		RemainingResends  uint32             `json:"remainingResends,omitempty"`
		ResendsIn         int64              `json:"resendsIn,omitempty"`
		ExpiresIn         int64              `json:"expiresIn"`
		Message           string             `json:"message,omitempty"`
		DebugInfo         string             `json:"debugInfo,omitempty"`
	}

	// ErrorConfirmOperationResponse - .
	ErrorConfirmOperationResponse struct {
		OperationStatus ConfirmOperationStatus  `json:"operationStatus,omitempty"`
		Errors          []mrresp.ErrorAttribute `json:"errors"`
	}

	// ConfirmOperationStatus - информация об оставшихся попытках и времени действия операции.
	// Поля RemainingResends и ResendsIn не используются для пароля и TOTP.
	ConfirmOperationStatus struct {
		RemainingAttempts uint32 `json:"remainingAttempts"`
		RemainingResends  uint32 `json:"remainingResends"`
		ResendsIn         int64  `json:"resendsIn"`
		ExpiresIn         int64  `json:"expiresIn"`
		DebugInfo         string `json:"debugInfo,omitempty"`
	}
)

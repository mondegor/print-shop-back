package model

type (
	// ChangeEmailRequest - запрос на изменение емаила пользователя.
	ChangeEmailRequest struct {
		NewEmail string `json:"newEmail" validate:"required,min=7,max=64,tag_email"`
	}

	// ChangePhoneRequest - запрос на установку/изменение телефона пользователя.
	ChangePhoneRequest struct {
		NewPhone string `json:"newPhone" validate:"required,min=10,max=32,tag_phone"`
	}

	// ChangePasswordRequest - запрос на установку/изменение пароля пользователя (2FA).
	ChangePasswordRequest struct {
		NewPassword string `json:"newPassword" validate:"required,min=8,max=32,tag_password"`
	}

	// ApplyOperationRequest - запрос на подтверждение операции.
	ApplyOperationRequest = OperationTokenRequest

	// GeneratedPasswordResponse - информация о надёжности пароля.
	GeneratedPasswordResponse struct {
		Password string `json:"password"`
	}

	// TOTPGeneratorResponse - .
	TOTPGeneratorResponse struct {
		URL     string `json:"totpUrl"`
		Message string `json:"message,omitempty"`
	}
)

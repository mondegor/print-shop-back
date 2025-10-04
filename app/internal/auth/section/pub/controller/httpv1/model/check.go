package model

type (
	// CheckLoginRequest - запрос на проверку свободен ли указанный емаил/телефон.
	CheckLoginRequest = AuthorizeUserRequest

	// CheckPasswordStrengthRequest - запрос на проверку надёжности указанного пароля.
	CheckPasswordStrengthRequest struct {
		Password string `json:"password" validate:"required,min=8,max=32,tag_password"`
	}

	// CheckPasswordStrengthResponse - информация о надёжности пароля.
	CheckPasswordStrengthResponse struct {
		Strength string `json:"strength"`
	}
)

package model

import (
	"github.com/mondegor/go-components/mrauth/enum"
)

type (
	// CreateUserRequest - запрос на создание нового пользователя.
	CreateUserRequest struct {
		Realm     string `json:"realm" validate:"required,max=32,tag_name"`
		UserEmail string `json:"userEmail" validate:"required,min=7,max=64,tag_email"`
	}

	// AuthorizeUserRequest - запрос на авторизацию пользователя в системе.
	AuthorizeUserRequest struct {
		Realm     string `json:"realm" validate:"required,max=32,tag_name"`
		UserLogin string `json:"userLogin" validate:"required,min=7,max=64,tag_email_phone"`
	}

	// LoginByTokenRequest - запрос на авторизацию пользователя в системе.
	LoginByTokenRequest struct {
		Token     string `json:"token" validate:"required,min=64,max=128"`
		Secret    string `json:"secret,omitempty" validate:"omitempty,min=4,max=32"`
		UseCookie bool   `json:"useCookie,omitempty"` // OPTIONAL
	}

	// ContinueSessionRequest - запрос на подтверждение операции.
	ContinueSessionRequest struct {
		RefreshToken string `json:"refreshToken" validate:"required,min=64,max=128"`
	}

	// SuccessAccessResponse - запрос на авторизацию пользователя в системе.
	SuccessAccessResponse struct {
		AccessToken  string `json:"accessToken"`
		ExpiresIn    uint32 `json:"expiresIn"`
		RefreshToken string `json:"refreshToken,omitempty"` // can be in cookie
		Message      string `json:"message,omitempty"`      // OPTIONAL
	}

	// UserInfoResponse - запрос на авторизацию пользователя в системе.
	UserInfoResponse struct {
		Email        string           `json:"email"`
		Phone        string           `json:"phone,omitempty"`
		LangCode     string           `json:"lang"`
		LastLoginIP  string           `json:"lastLoginIp"`
		LastLoggedAt string           `json:"lastLoggedAt"`
		Auth2faType  enum.Auth2faType `json:"auth2faType"`
		Realms       []UserRealm      `json:"realms"`
		Status       enum.UserStatus  `json:"status"`
		// CreatedAt    time.Time       `json:"createdAt"`
		// UpdatedAt    time.Time       `json:"updatedAt"`
	}

	// UserRealm - запрос на авторизацию пользователя в системе.
	UserRealm struct {
		Name     string `json:"name"`
		UserKind string `json:"userKind"`
		// CreatedAt time.Time `json:"createdAt"`
		// UpdatedAt time.Time `json:"updatedAt"`
	}
)

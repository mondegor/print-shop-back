package auth

import (
	"time"
)

type (
	// UserRealm - comment struct.
	UserRealm struct {
		Name             string           `yaml:"name"`
		AuthToken        Token            `yaml:"auth_token"`
		UserKinds        []UserKind       `yaml:"user_kinds"`
		RegisterUserKind string           `yaml:"register_user_kind"`
		OperationConfirm OperationConfirm `yaml:"operation_confirm"`
	}

	// Token - comment struct.
	Token struct {
		AccessType    string        `yaml:"access_type"`
		AccessExpiry  time.Duration `yaml:"access_expiry"`
		RefreshExpiry time.Duration `yaml:"refresh_expiry"`
		Length        uint32        `yaml:"length"`
	}

	// UserKind - comment struct.
	UserKind struct {
		Name       string   `yaml:"name"`
		Roles      []string `yaml:"roles"`
		SessionMax uint32   `yaml:"session_max"`
	}

	// OperationConfirm - comment struct.
	OperationConfirm struct {
		TokenLength   uint32        `yaml:"token_length"`
		CodeLength    uint32        `yaml:"code_length"`
		SessionExpiry time.Duration `yaml:"session_expiry"`
		SendByEmail   CodeSender    `yaml:"send_by_email"`
		SendByPhone   CodeSender    `yaml:"send_by_phone"`
	}

	// CodeSender - comment struct.
	CodeSender struct {
		MaxAttempts   uint32        `yaml:"max_attempts"`
		MaxResends    uint32        `yaml:"max_resends"`
		MinResendTime time.Duration `yaml:"min_resend_time"`
	}

	// JWTConfig - comment struct.
	JWTConfig struct {
		Method string
		Secret []byte
	}
)

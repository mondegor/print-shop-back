package auth

import (
	"time"

	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/pkg/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Logger              mrlog.Logger
		EventEmitter        mrevent.Emitter
		UseCaseErrorWrapper mrerr.UseCaseErrorWrapper
		StorageErrorWrapper mrerr.ErrorWrapper
		DBConnManager       mrstorage.DBConnManager
		RequestParsers      RequestParsers
		ResponseSender      mrserver.FileResponseSender
		NotifierAPI         mrnotifier.NoticeProducer
		Locker              mrlock.Locker
		UserRealms          []UserRealm
		OperationConfirm    OperationConfirm
		JWT                 JWT
		WithDebugInfo       bool
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}

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

	// JWT - comment struct.
	JWT struct {
		Method string
		Secret []byte
	}
)

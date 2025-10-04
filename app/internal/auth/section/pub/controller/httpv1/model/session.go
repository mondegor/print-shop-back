package model

type (
	// CloseSessionsRequest - запрос на подтверждение операции.
	CloseSessionsRequest struct {
		Hashes []string `json:"hashes" validate:"required,gte=1,dive,max=16,tag_session_hash"`
	}
)

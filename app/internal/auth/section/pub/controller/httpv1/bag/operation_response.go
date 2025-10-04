package bag

import (
	"github.com/mondegor/go-components/mrauth/entity"
	"github.com/mondegor/go-components/mrauth/enum"
	"github.com/mondegor/go-sysmess/mrlib/exttime"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/model"
)

// OperationResponse - comment struct.
type OperationResponse struct {
	withDebugInfo bool
}

// NewOperationResponse - создаёт объект OperationResponse.
func NewOperationResponse(withDebugInfo bool) *OperationResponse {
	return &OperationResponse{
		withDebugInfo: withDebugInfo,
	}
}

// NewConfirmOperation - comment method.
func (r *OperationResponse) NewConfirmOperation(operation entity.SecureOperation, message string) model.WaitingConfirmOperationResponse {
	return model.WaitingConfirmOperationResponse{
		Token:             operation.Token,
		ConfirmMethod:     r.operationAction(&operation).Method,
		RemainingAttempts: operation.RemainingAttempts,
		RemainingResends:  operation.RemainingResends,
		ResendsIn:         exttime.TimeLeftInSec(operation.ResendsAt),
		ExpiresIn:         exttime.TimeLeftInSec(operation.ExpiresAt),
		Message:           message,
		DebugInfo:         r.debugInfo(&operation),
	}
}

// NewErrorConfirmOperation - comment method.
func (r *OperationResponse) NewErrorConfirmOperation(operation entity.SecureOperation, lz mrcore.Localizer, err error) model.ErrorConfirmOperationResponse {
	return model.ErrorConfirmOperationResponse{
		OperationStatus: r.newOperationStatus(&operation),
		Errors: []mrresp.ErrorAttribute{
			mrresp.NewErrorAttribute(lz, err, r.withDebugInfo),
		},
	}
}

func (r *OperationResponse) newOperationStatus(operation *entity.SecureOperation) model.ConfirmOperationStatus {
	return model.ConfirmOperationStatus{
		RemainingAttempts: operation.RemainingAttempts,
		RemainingResends:  operation.RemainingResends,
		ResendsIn:         exttime.TimeLeftInSec(operation.ResendsAt),
		ExpiresIn:         exttime.TimeLeftInSec(operation.ExpiresAt),
		DebugInfo:         r.debugInfo(operation),
	}
}

func (r *OperationResponse) operationAction(operation *entity.SecureOperation) entity.ConfirmAction {
	for i := range operation.Actions {
		if !operation.Actions[i].Confirmed {
			return operation.Actions[i]
		}
	}

	return entity.ConfirmAction{}
}

func (r *OperationResponse) debugInfo(operation *entity.SecureOperation) string {
	if !r.withDebugInfo {
		return ""
	}

	action := r.operationAction(operation)

	if action.Method == enum.ConfirmMethodEmail || action.Method == enum.ConfirmMethodPhone {
		return "Confirm code: " + action.Secret
	}

	return ""
}

package mrhttp

import "encoding/json"

const (
    AppErrorAttributeNameSystem = "system"
)

type (
    // application/problem+json:
    AppErrorResponse struct {
        Title string `json:"title"`
        Details string `json:"details"`
        Request string `json:"request"`
        Time string `json:"time"`
        ErrorTraceId string `json:"errorTraceId,omitempty"`
    }

    // application/json:
    AppErrorListResponse []AppErrorAttribute

    AppErrorAttribute struct {
        Id string `json:"id"`
        Value string `json:"value"`
    }
)

func (ar *AppErrorResponse) Marshal() []byte {
    bytes, err := json.Marshal(ar)

    if err != nil {
        return nil
    }

    return bytes
}

func (a *AppErrorListResponse) Add(id string, value string) {
    *a = append(*a, AppErrorAttribute{Id: id, Value: value})
}

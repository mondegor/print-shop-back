package entity

type (
    ResourceStatusFlow map[ResourceStatus][]ResourceStatus
)

var (
    ResourceStatusFlowDefault = ResourceStatusFlow{
        ResourceStatusDraft: []ResourceStatus{
            ResourceStatusPublished,
        },
        ResourceStatusHidden: []ResourceStatus{
            ResourceStatusPublished,
        },
        ResourceStatusPublished: []ResourceStatus{
            ResourceStatusHidden,
        },
    }
)

func (isf ResourceStatusFlow) Check(statusFrom ResourceStatus, statusTo ResourceStatus) bool {
    for _, status := range isf.getPossibleToStatuses(statusFrom) {
        if statusTo == status {
            return true
        }
    }

    return false
}

// getPossibleToStatuses
// Возвращается список статусов в которые можно переключить указанный статус
func (isf ResourceStatusFlow) getPossibleToStatuses(statusFrom ResourceStatus) []ResourceStatus {
    if statuses, ok := isf[statusFrom]; ok {
        return statuses
    }

    return []ResourceStatus{}
}

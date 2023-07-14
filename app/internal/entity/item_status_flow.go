package entity

type ItemStatusFlow map[ItemStatus][]ItemStatus

var (
    ItemStatusFlowDefault = ItemStatusFlow{
        ItemStatusDraft: []ItemStatus{
            ItemStatusEnabled,
            ItemStatusDisabled,
            ItemStatusRemoved,
        },
        ItemStatusEnabled: []ItemStatus{
            ItemStatusDisabled,
            ItemStatusRemoved,
        },
        ItemStatusDisabled: []ItemStatus{
            ItemStatusEnabled,
            ItemStatusRemoved,
        },
        ItemStatusRemoved: []ItemStatus{},
    }

    ItemStatusFlowOnlyRemove = ItemStatusFlow{
        ItemStatusEnabled: []ItemStatus{
            ItemStatusRemoved,
        },
        ItemStatusRemoved: []ItemStatus{},
    }
)

func (isf ItemStatusFlow) Check(statusFrom ItemStatus, statusTo ItemStatus) bool {
    for _, status := range isf.getPossibleToStatuses(statusFrom) {
        if statusTo == status {
            return true
        }
    }

    return false
}

// getPossibleToStatuses
// Возвращается список статусов в которые можно переключить указанный статус
func (isf ItemStatusFlow) getPossibleToStatuses(statusFrom ItemStatus) []ItemStatus {
    if statuses, ok := isf[statusFrom]; ok {
        return statuses
    }

    return []ItemStatus{}
}

//// getPossibleFromStatuses
//// Возвращается список статусов из которых можно переключиться в указанный статус
//func (isf ItemStatusFlow) getPossibleFromStatuses(statusTo ItemStatus) []ItemStatus {
//  if _, ok := isf[statusTo]; !ok {
//      return nil
//  }
//
//  var statuses []ItemStatus
//
//  for statusFrom, statusesTo := range isf {
//      for _, status := range statusesTo {
//          if statusTo == status {
//              statuses = append(statuses, statusFrom)
//              break
//          }
//      }
//  }
//
//  return statuses
//}

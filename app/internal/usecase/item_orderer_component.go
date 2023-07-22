package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

const (
    // orderFieldStep mrentity.Int64 = 1024 * 1024
    orderFieldStep mrentity.Int64 = 1
)

type ItemOrderer struct {
    storage ItemOrdererStorage
}

func NewItemOrdererComponent(storage ItemOrdererStorage) *ItemOrderer {
    return &ItemOrderer{
        storage: storage,
    }
}

func (f *ItemOrderer) WithMetaData(meta ItemMetaData) ItemOrdererComponent {
    return &ItemOrderer{
        storage: f.storage.WithMetaData(meta),
    }
}

func (f *ItemOrderer) InsertToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    firstNode := entity.ItemOrdererNode{}
    err := f.storage.LoadFirstNode(ctx, &firstNode)

    if err != nil {
        return err
    }

    if firstNode.Id == nodeId {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    err = f.storage.UpdateNodePrevId(ctx, firstNode.Id, mrentity.ZeronullInt32(nodeId))

    if err != nil {
        return err
    }

    currentNode := entity.ItemOrdererNode{
        Id: nodeId,
        PrevId: 0,
        NextId: mrentity.ZeronullInt32(firstNode.Id),
        OrderField: firstNode.OrderField / 2,
    }

    if currentNode.OrderField < 1 {
        err = f.storage.RecalcOrderField(ctx, 0, 2 * orderFieldStep)

        if err != nil {
            return err
        }

        currentNode.OrderField = mrentity.ZeronullInt64(orderFieldStep)
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::InsertToFirst: id=%d", entity.ModelNameItemOrderer, nodeId)

    return nil
}

func (f *ItemOrderer) InsertToLast(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    lastNode := entity.ItemOrdererNode{}
    err := f.storage.LoadLastNode(ctx, &lastNode)

    if err != nil {
        return err
    }

    if lastNode.Id == nodeId {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    err = f.storage.UpdateNodeNextId(ctx, lastNode.Id, mrentity.ZeronullInt32(nodeId))

    if err != nil {
        return err
    }

    currentNode := entity.ItemOrdererNode{
        Id: nodeId,
        PrevId: mrentity.ZeronullInt32(lastNode.Id),
        NextId: 0,
        OrderField: lastNode.OrderField + mrentity.ZeronullInt64(orderFieldStep),
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::InsertToLast: id=%d", entity.ModelNameItemOrderer, nodeId)

    return nil
}

func (f *ItemOrderer) MoveToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    currentNode := entity.ItemOrdererNode{Id: nodeId}

    firstNode := entity.ItemOrdererNode{}
    err := f.storage.LoadFirstNode(ctx, &firstNode)

    if err != nil {
        return err
    }

    if firstNode.Id == currentNode.Id {
        if firstNode.OrderField == 0 {
            currentNode.OrderField = mrentity.ZeronullInt64(orderFieldStep)
            err = f.storage.UpdateNode(ctx, &currentNode)

            if err != nil {
                return err
            }
        }

        return nil
    }

    err = f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    if mrentity.KeyInt32(currentNode.NextId) == firstNode.Id {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"currentNode.Id": currentNode.Id, "currentNode.NextId": currentNode.NextId})
    }

    err = f.storage.UpdateNodePrevId(ctx, firstNode.Id, mrentity.ZeronullInt32(currentNode.Id))

    if err != nil {
        return err
    }

    if currentNode.PrevId > 0 {
        err = f.storage.UpdateNodeNextId(ctx, mrentity.KeyInt32(currentNode.PrevId), currentNode.NextId)

        if err != nil {
            return err
        }
    }

    if currentNode.NextId > 0 {
        err = f.storage.UpdateNodePrevId(ctx, mrentity.KeyInt32(currentNode.NextId), currentNode.PrevId)

        if err != nil {
            return err
        }
    }

    currentNode.PrevId = 0
    currentNode.NextId = mrentity.ZeronullInt32(firstNode.Id)
    currentNode.OrderField = firstNode.OrderField / 2

    if currentNode.OrderField < 1 {
        err = f.storage.RecalcOrderField(ctx, 0, 2 * orderFieldStep)

        if err != nil {
            return err
        }

        currentNode.OrderField = mrentity.ZeronullInt64(orderFieldStep)
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::MoveToFirst: id=%d", entity.ModelNameItemOrderer, nodeId)

    return nil
}

func (f *ItemOrderer) MoveToLast(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    currentNode := entity.ItemOrdererNode{Id: nodeId}

    lastNode := entity.ItemOrdererNode{}
    err := f.storage.LoadLastNode(ctx, &lastNode)

    if err != nil {
        return err
    }

    if lastNode.Id == currentNode.Id {
        if lastNode.OrderField == 0 {
            currentNode.OrderField = mrentity.ZeronullInt64(orderFieldStep)
            err = f.storage.UpdateNode(ctx, &currentNode)

            if err != nil {
                return err
            }
        }

        return nil
    }

    err = f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    if mrentity.KeyInt32(currentNode.PrevId) == lastNode.Id {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"currentNode.Id": currentNode.Id, "currentNode.PrevId": currentNode.PrevId})
    }

    err = f.storage.UpdateNodeNextId(ctx, lastNode.Id, mrentity.ZeronullInt32(currentNode.Id))

    if err != nil {
        return err
    }

    if currentNode.PrevId > 0 {
        err = f.storage.UpdateNodeNextId(ctx, mrentity.KeyInt32(currentNode.PrevId), currentNode.NextId)

        if err != nil {
            return err
        }
    }

    if currentNode.NextId > 0 {
        err = f.storage.UpdateNodePrevId(ctx, mrentity.KeyInt32(currentNode.NextId), currentNode.PrevId)

        if err != nil {
            return err
        }
    }

    currentNode.PrevId = mrentity.ZeronullInt32(lastNode.Id)
    currentNode.NextId = 0
    currentNode.OrderField = lastNode.OrderField + mrentity.ZeronullInt64(orderFieldStep)

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::MoveToLast: id=%d", entity.ModelNameItemOrderer, nodeId)

    return nil
}

func (f *ItemOrderer) MoveAfterId(ctx context.Context, nodeId mrentity.KeyInt32, afterNodeId mrentity.KeyInt32) error {
    if afterNodeId < 1 {
        return f.MoveToFirst(ctx, nodeId)
    }

    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId})
    }

    if nodeId == afterNodeId {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"nodeId": nodeId, "afterNodeId": afterNodeId})
    }

    currentNode := entity.ItemOrdererNode{Id: nodeId}
    err := f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    if mrentity.KeyInt32(currentNode.PrevId) == afterNodeId {
        return nil
    }

    afterNode := entity.ItemOrdererNode{Id: afterNodeId}
    err = f.storage.LoadNode(ctx, &afterNode)

    if err != nil {
        return err
    }

    afterNextNode := entity.ItemOrdererNode{Id: mrentity.KeyInt32(afterNode.NextId)}

    if afterNextNode.Id > 0 {
        err = f.storage.LoadNode(ctx, &afterNextNode)

        if err != nil {
            return err
        }
    }

    err = f.storage.UpdateNodeNextId(ctx, afterNode.Id, mrentity.ZeronullInt32(currentNode.Id))

    if err != nil {
        return err
    }

    if afterNextNode.Id > 0 {
        err = f.storage.UpdateNodePrevId(ctx, afterNextNode.Id, mrentity.ZeronullInt32(currentNode.Id))

        if err != nil {
            return err
        }
    }

    if currentNode.PrevId > 0 {
        err = f.storage.UpdateNodeNextId(ctx, mrentity.KeyInt32(currentNode.PrevId), currentNode.NextId)

        if err != nil {
            return err
        }
    }

    if currentNode.NextId > 0 {
        err = f.storage.UpdateNodePrevId(ctx, mrentity.KeyInt32(currentNode.NextId), currentNode.PrevId)

        if err != nil {
            return err
        }
    }

    currentNode.PrevId = mrentity.ZeronullInt32(afterNode.Id)
    currentNode.NextId = mrentity.ZeronullInt32(afterNextNode.Id)
    currentNode.OrderField = (afterNode.OrderField + afterNextNode.OrderField) / 2

    if currentNode.OrderField <= afterNode.OrderField {
        if afterNextNode.Id > 0 {
            err = f.storage.RecalcOrderField(ctx, mrentity.Int64(afterNode.OrderField), 2 * orderFieldStep)

            if err != nil {
                return err
            }
        }

        currentNode.OrderField = afterNode.OrderField + mrentity.ZeronullInt64(orderFieldStep)
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::MoveAfterId: id=%d, afterId=%d, ", entity.ModelNameItemOrderer, nodeId, afterNodeId)

    return nil
}

func (f *ItemOrderer) Unlink(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return f.MoveToFirst(ctx, nodeId)
    }

    currentNode := entity.ItemOrdererNode{Id: nodeId}
    err := f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    if currentNode.PrevId == 0 &&
        currentNode.NextId == 0 &&
        currentNode.OrderField == 0 {
        return nil
    }

    if currentNode.PrevId > 0 {
        err = f.storage.UpdateNodeNextId(ctx, mrentity.KeyInt32(currentNode.PrevId), currentNode.NextId)

        if err != nil {
            return err
        }
    }

    if currentNode.NextId > 0 {
        err = f.storage.UpdateNodePrevId(ctx, mrentity.KeyInt32(currentNode.NextId), currentNode.PrevId)

        if err != nil {
            return err
        }
    }

    currentNode.PrevId = 0
    currentNode.NextId = 0
    currentNode.OrderField = 0

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    f.logger(ctx).Event("%s::Unlink: id=%d", entity.ModelNameItemOrderer, nodeId)

    return nil
}

func (f *ItemOrderer) logger(ctx context.Context) mrapp.Logger {
   return mrcontext.GetLogger(ctx)
}

package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "context"
)

const (
    // orderFieldStep mrentity.ZeronullInt64 = 1024 * 1024
    orderFieldStep mrentity.ZeronullInt64 = 1024
)

type FormFieldItemOrderer struct {
    storage FormFieldItemOrdererStorage
    errorHelper *mrerr.Helper
}

func NewFormFieldItemOrderer(storage FormFieldItemOrdererStorage, errorHelper *mrerr.Helper) *FormFieldItemOrderer {
    return &FormFieldItemOrderer{
        storage: storage,
        errorHelper: errorHelper,
    }
}

func (f *FormFieldItemOrderer) InsertToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    firstNode := entity.Node{}
    err := f.storage.LoadFirstNode(ctx, &firstNode)

    if err != nil {
        return err
    }

    if firstNode.Id == nodeId {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    err = f.storage.UpdateNodePrevId(ctx, firstNode.Id, mrentity.ZeronullInt32(nodeId))

    if err != nil {
        return err
    }

    currentNode := entity.Node{
        Id: nodeId,
        PrevId: 0,
        NextId: mrentity.ZeronullInt32(firstNode.Id),
        OrderField: firstNode.OrderField / 2,
    }

    if currentNode.OrderField < 1 {
        err = f.storage.RecalcOrderField(ctx, 0, 2 * orderFieldStep)
        currentNode.OrderField = orderFieldStep
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    return nil
}

func (f *FormFieldItemOrderer) InsertToLast(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    lastNode := entity.Node{}
    err := f.storage.LoadLastNode(ctx, &lastNode)

    if err != nil {
        return err
    }

    if lastNode.Id == nodeId {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    err = f.storage.UpdateNodeNextId(ctx, lastNode.Id, mrentity.ZeronullInt32(nodeId))

    if err != nil {
        return err
    }

    currentNode := entity.Node{
        Id: nodeId,
        PrevId: mrentity.ZeronullInt32(lastNode.Id),
        NextId: 0,
        OrderField: lastNode.OrderField + orderFieldStep,
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    return nil
}

func (f *FormFieldItemOrderer) MoveToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    currentNode := entity.Node{Id: nodeId}

    firstNode := entity.Node{}
    err := f.storage.LoadFirstNode(ctx, &firstNode)

    if err != nil {
        return err
    }

    if firstNode.Id == currentNode.Id {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("currentNode.Id=%d", currentNode.Id)
    }

    err = f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
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
        currentNode.OrderField = orderFieldStep
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    return nil
}

func (f *FormFieldItemOrderer) MoveToLast(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    currentNode := entity.Node{Id: nodeId}

    lastNode := entity.Node{}
    err := f.storage.LoadLastNode(ctx, &lastNode)

    if err != nil {
        return err
    }

    if lastNode.Id == currentNode.Id {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("currentNode.Id=%d", currentNode.Id)
    }

    err = f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
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
    currentNode.OrderField = lastNode.OrderField + orderFieldStep

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    return nil
}

func (f *FormFieldItemOrderer) MoveAfterId(ctx context.Context, nodeId mrentity.KeyInt32, afterNodeId mrentity.KeyInt32) error {
    if afterNodeId < 1 {
        return f.MoveToFirst(ctx, nodeId)
    }

    if nodeId < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("nodeId=%d", nodeId)
    }

    currentNode := entity.Node{Id: nodeId}
    err := f.storage.LoadNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    afterNode := entity.Node{Id: afterNodeId}
    err = f.storage.LoadNode(ctx, &afterNode)

    if err != nil {
        return err
    }

    afterNextNode := entity.Node{Id: mrentity.KeyInt32(afterNode.NextId)}

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

    if currentNode.OrderField < afterNode.OrderField {
        if afterNextNode.Id > 0 {
            err = f.storage.RecalcOrderField(ctx, afterNode.OrderField, 2 * orderFieldStep)
        }

        currentNode.OrderField = afterNode.OrderField + orderFieldStep
    }

    err = f.storage.UpdateNode(ctx, &currentNode)

    if err != nil {
        return err
    }

    return nil
}

func (f *FormFieldItemOrderer) Unlink(ctx context.Context, nodeId mrentity.KeyInt32) error {
    if nodeId < 1 {
        return f.MoveToFirst(ctx, nodeId)
    }

    currentNode := entity.Node{Id: nodeId}
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

    return nil
}

//func (f *FormFieldItemOrderer) logger(ctx context.Context) mrapp.Logger {
//    return mrcontext.GetLogger(ctx)
//}

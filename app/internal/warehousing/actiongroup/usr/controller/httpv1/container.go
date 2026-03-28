package httpv1

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
	"github.com/mondegor/print-shop-back/internal/warehousing/module"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	containerListURL      = "/v1/warehousing/containers"
	containerGroupListURL = "/v1/warehousing/container-groups"
	containerTagsURL      = "/v1/warehousing/containers/tags"
	containerImagesURL    = "/v1/warehousing/containers/images"
)

type (
	// Container - comment struct.
	Container struct {
		parser                 validate.RequestParser
		sender                 mrserver.ResponseSender
		serviceContainer       usr.ContainerService
		useCaseCreateContainer createContainerUseCase
	}

	createContainerUseCase interface {
		Execute(ctx context.Context, item dto.CreateStockContainer) (dto.CreateStockContainerResult, error)
	}
)

// NewContainer - создаёт контроллер Container.
func NewContainer(
	parser validate.RequestParser,
	sender mrserver.ResponseSender,
	serviceContainer usr.ContainerService,
	createContainerUseCase createContainerUseCase,
) *Container {
	return &Container{
		parser:                 parser,
		sender:                 sender,
		serviceContainer:       serviceContainer,
		useCaseCreateContainer: createContainerUseCase,
	}
}

// Handlers - возвращает обработчики контроллера Container.
func (ht *Container) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: containerListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: containerListURL, Func: ht.CreateContainer},
		{Method: http.MethodPost, URL: containerGroupListURL, Func: ht.CreateGroup},
		{Method: http.MethodPut, URL: containerTagsURL, Func: ht.SaveTags},
		{Method: http.MethodPost, URL: containerImagesURL, Func: ht.AddImages},
		{Method: http.MethodDelete, URL: containerImagesURL, Func: ht.DeleteImages},
	}
}

// GetList - comment method.
func (ht *Container) GetList(w http.ResponseWriter, r *http.Request) error {
	items, hasNext, err := ht.serviceContainer.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ContainerListResponse{
			Items:   items,
			Cursor:  entity.CreateContainerCursorValue(items),
			HasNext: hasNext,
		},
	)
}

func (ht *Container) listParams(r *http.Request) dto.ContainerParams {
	return dto.ContainerParams{
		AccountID: ht.parser.UserID(r),
		Filter: dto.ContainerFilter{
			SearchCode: ht.parser.FilterString(r, module.ParamNameFilterSearchContainerCode),
			SearchTags: ht.parser.FilterStringList(r, module.ParamNameFilterSearchContainerTags),
		},
		Cursor: xtype.NewContainerCursor(ht.parser.CursorParams(r)),
	}
}

// CreateContainer - comment method.
func (ht *Container) CreateContainer(w http.ResponseWriter, r *http.Request) error {
	req := CreateContainerRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := dto.CreateStockContainer{
		AccountID:  ht.parser.UserID(r),
		Kind:       locationkind.Container,
		Code:       req.Code,
		Volume:     req.Volume,
		Tags:       req.Tags,
		Images:     req.Images,
		LocationID: req.LocationID,
		Quantity:   req.ExemplarQuantity,
	}

	createdContainer, err := ht.useCaseCreateContainer.Execute(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessCreatedContainerResponse{
			ContainerID: createdContainer.ID,
			Code:        createdContainer.Code,
			Marker:      createdContainer.Marker,
			StockID:     createdContainer.StockID,
		},
	)
}

// CreateGroup - comment method.
func (ht *Container) CreateGroup(w http.ResponseWriter, r *http.Request) error {
	req := CreateGroupRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := dto.CreateStockContainer{
		AccountID:  ht.parser.UserID(r),
		Kind:       locationkind.Group,
		Code:       req.Code,
		Volume:     req.Volume,
		LocationID: req.LocationID,
		Quantity:   1, // group has only one
	}

	createdContainer, err := ht.useCaseCreateContainer.Execute(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessCreatedContainerResponse{
			ContainerID: createdContainer.ID,
			Code:        createdContainer.Code,
			Marker:      createdContainer.Marker,
			StockID:     createdContainer.StockID,
		},
	)
}

// SaveTags - comment method.
func (ht *Container) SaveTags(w http.ResponseWriter, r *http.Request) error {
	req := SaveTagsRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	tagVersion, err := ht.serviceContainer.SaveTags(
		r.Context(),
		entity.UpdateContainerTags{
			ID:         req.ContainerID,
			AccountID:  ht.parser.UserID(r),
			TagVersion: req.TagVersion,
			Tags:       req.Tags,
		},
	)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		TagVersionResponse{
			TagVersion: tagVersion,
		},
	)
}

// AddImages - comment method.
func (ht *Container) AddImages(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// DeleteImages - comment method.
func (ht *Container) DeleteImages(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (ht *Container) getItemValue(r *http.Request) string {
	return ht.parser.PathParamString(r, "value")
}

func (ht *Container) wrapError(err error, r *http.Request) error {
	if errors.Is(err, errors.ErrRecordNotFound) {
		return module.ErrContainerNotFound.Wrap(err, ht.getItemValue(r))
	}

	return err
}

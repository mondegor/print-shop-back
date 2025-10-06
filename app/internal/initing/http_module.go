package initing

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// HttpModule - информация для создания и инициализации всех контроллеров Http модуля.
	HttpModule struct {
		Name       string
		Permission string

		// InitSharedComponents - функция вызывается перед созданием первого контроллера,
		// с помощью нее можно инициализировать общие переменные доступные
		// всем контроллерам, которые можно только получить предварительно обработав ошибки.
		InitSharedComponents func() (err error)
		Controllers          []HttpController
	}

	// HttpController - информация для создания и инициализации Http контроллера, и всех его обработчиков.
	HttpController struct {
		Name       string
		Permission string
		Create     func() (mrserver.HttpController, error)
	}
)

// CreateHttpControllers - создание инициализация всех контроллеров для указанных модулей.
func CreateHttpControllers(logger mrlog.Logger, modules []HttpModule, operations ...mrinit.PrepareHandlerFunc) (list []mrserver.HttpController, err error) {
	for _, module := range modules {
		mrlog.Info(logger, "Create and init module", "module", module.Name, "permission", module.Permission)

		if module.InitSharedComponents != nil {
			if err := module.InitSharedComponents(); err != nil {
				return nil, fmt.Errorf("init shared components: %w", err)
			}
		}

		for _, c := range module.Controllers {
			if c.Create == nil {
				return nil, fmt.Errorf("create controller for module %s (c.Create = nil)", module.Name)
			}

			controller, err := c.Create()
			if err != nil {
				return nil, err
			}

			if c.Name != "" || c.Permission != "" {
				mrlog.Info(logger, "Create and init controller", "controller", c.Name, "permission", c.Permission)
			}

			// если разрешение контроллера не указано, то используется разрешение его модуля.
			if c.Permission == "" {
				c.Permission = module.Permission
			}

			list = append(
				list,
				mrinit.PrepareController(
					controller,
					append([]mrinit.PrepareHandlerFunc{mrinit.WithPermission(c.Permission)}, operations...)...,
				),
			)
		}
	}

	return list, nil
}

package rest

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mraccess/section"

	"github.com/mondegor/print-shop-back/internal/app"
)

// initRoutingSections - создаёт объект section.RoutingSection с указанными настройками.
func initRoutingSections(opts app.Options) (name2section map[string]*section.RoutingSection) {
	name2section = make(map[string]*section.RoutingSection, len(opts.Cfg.AccessControl.RoutingSections))

	for _, sect := range opts.Cfg.AccessControl.RoutingSections {
		mrlog.Info(
			opts.Logger,
			fmt.Sprintf(
				"Init section '%s' with root path '%s' and privilege '%s'",
				sect.Name, sect.BasePath, sect.Privilege,
			),
		)

		name2section[sect.Name] = section.NewRoutingSection(
			sect.Name,
			sect.BasePath,
			sect.Privilege,
		)
	}

	return name2section
}

package initing

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mraccess/section"
)

type (
	// RoutingSection - comment struct.
	RoutingSection struct {
		Name      string `yaml:"name"`
		BasePath  string `yaml:"base_path"`
		Privilege string `yaml:"privilege"`
	}
)

// InitRoutingSections - создаёт объект section.RoutingSection с указанными настройками.
func InitRoutingSections(logger mrlog.Logger, sections []RoutingSection) (name2section map[string]*section.RoutingSection) {
	name2section = make(map[string]*section.RoutingSection, len(sections))

	for _, sect := range sections {
		mrlog.Info(
			logger,
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

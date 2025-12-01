package factory

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrpath"

	"github.com/mondegor/print-shop-back/config"
)

// InitImageURLBuilder - создаёт объект placeholderpath.Builder.
func InitImageURLBuilder(cfg config.Config) (mrpath.Builder, error) {
	return mrpath.NewPlaceholder(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/")+
			"/"+
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.BasePath, "/"),
		mrpath.Placeholder,
	)
}

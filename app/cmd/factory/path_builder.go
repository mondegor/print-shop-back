package factory

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrpath"

	"print-shop-back/config"
)

// InitImageURLBuilder - создаёт объект placeholderpath.Builder.
func InitImageURLBuilder(cfg config.Config) (mrpath.Builder, error) {
	return mrpath.NewPlaceholder(
		strings.TrimRight(cfg.ModuleSettings.FileStation.ImageProxyHost, "/")+
			"/"+
			strings.TrimLeft(cfg.ModuleSettings.FileStation.ImageProxyBasePath, "/"),
		mrpath.Placeholder,
	)
}

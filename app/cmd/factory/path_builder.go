package factory

import (
	"strings"

	"github.com/mondegor/go-webcore/mrpath/placeholderpath"

	"github.com/mondegor/print-shop-back/config"
)

// InitImageURLBuilder - создаёт объект placeholderpath.Builder.
func InitImageURLBuilder(cfg config.Config) (*placeholderpath.Builder, error) {
	return placeholderpath.New(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/")+
			"/"+
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.BasePath, "/"),
		placeholderpath.Placeholder,
	)
}

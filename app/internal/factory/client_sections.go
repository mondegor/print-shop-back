package factory

import (
	"print-shop-back/config"
	"print-shop-back/internal/global"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewClientSectionAdminAPI(cfg *config.Config, logger mrcore.Logger, access *mrperms.ModulesAccess) *mrperms.ClientSection {
	return newClientSection(
		mrperms.ClientSectionOptions{
			Caption:      global.SectionAdminAPICaption,
			RootPath:     global.SectionAdminAPIRootPath,
			Privilege:    cfg.ClientSections.AdminAPI.Privilege,
			AuthSecret:   cfg.ClientSections.AdminAPI.Auth.Secret,
			AuthAudience: cfg.ClientSections.AdminAPI.Auth.Audience,
			Access:       access,
		},
		logger,
	)
}

func NewClientSectionProviderAccountAPI(cfg *config.Config, logger mrcore.Logger, access *mrperms.ModulesAccess) *mrperms.ClientSection {
	return newClientSection(
		mrperms.ClientSectionOptions{
			Caption:      global.SectionProviderAccountAPICaption,
			RootPath:     global.SectionProviderAccountAPIRootPath,
			Privilege:    cfg.ClientSections.ProviderAccountAPI.Privilege,
			AuthSecret:   cfg.ClientSections.ProviderAccountAPI.Auth.Secret,
			AuthAudience: cfg.ClientSections.ProviderAccountAPI.Auth.Audience,
			Access:       access,
		},
		logger,
	)
}

func NewClientSectionPublicAPI(cfg *config.Config, logger mrcore.Logger, access *mrperms.ModulesAccess) *mrperms.ClientSection {
	return newClientSection(
		mrperms.ClientSectionOptions{
			Caption:      global.SectionPublicAPICaption,
			RootPath:     global.SectionPublicAPIRootPath,
			Privilege:    cfg.ClientSections.PublicAPI.Privilege,
			AuthSecret:   cfg.ClientSections.PublicAPI.Auth.Secret,
			AuthAudience: cfg.ClientSections.PublicAPI.Auth.Audience,
			Access:       access,
		},
		logger,
	)
}

func newClientSection(opt mrperms.ClientSectionOptions, logger mrcore.Logger) *mrperms.ClientSection {
	logger.Info("Init section %s with root path '%s' and privilege '%s'", opt.Caption, opt.RootPath, opt.Privilege)
	logger.Debug("secret=%s, audience: %s", opt.AuthSecret, opt.AuthAudience)

	return mrperms.NewClientSection(opt)
}

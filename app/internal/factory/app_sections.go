package factory

import (
	"print-shop-back/internal/modules"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrperms"
)

const (
	sectionAdminAPICaption  = "Admin API"
	sectionAdminAPIRootPath = "/adm/"

	sectionProviderAccountAPICaption  = "Provider Account API"
	sectionProviderAccountAPIRootPath = "/pacc/"

	sectionPublicAPICaption  = "Public API"
	sectionPublicAPIRootPath = "/"
)

func NewAppSectionAdminAPI(opts *modules.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		mrperms.AppSectionOptions{
			Caption:      sectionAdminAPICaption,
			RootPath:     sectionAdminAPIRootPath,
			Privilege:    opts.Cfg.AppSections.AdminAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.AdminAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.AdminAPI.Auth.Audience,
		},
		opts.AccessControl,
		opts.Logger,
	)
}

func NewAppSectionProviderAccountAPI(opts *modules.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		mrperms.AppSectionOptions{
			Caption:      sectionProviderAccountAPICaption,
			RootPath:     sectionProviderAccountAPIRootPath,
			Privilege:    opts.Cfg.AppSections.ProviderAccountAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.ProviderAccountAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.ProviderAccountAPI.Auth.Audience,
		},
		opts.AccessControl,
		opts.Logger,
	)
}

func NewAppSectionPublicAPI(opts *modules.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		mrperms.AppSectionOptions{
			Caption:      sectionPublicAPICaption,
			RootPath:     sectionPublicAPIRootPath,
			Privilege:    opts.Cfg.AppSections.PublicAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.PublicAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.PublicAPI.Auth.Audience,
		},
		opts.AccessControl,
		opts.Logger,
	)
}

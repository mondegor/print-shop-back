package factory

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrperms"
)

const (
	sectionAdminAPICaption  = "Admin API"
	sectionAdminAPIBasePath = "/adm/"

	sectionProvidersAPICaption  = "Providers API"
	sectionProvidersAPIBasePath = "/prov/"

	sectionPublicAPICaption  = "Public API"
	sectionPublicAPIBasePath = "/"
)

// NewAppSectionAdminAPI - создаёт объект mrperms.AppSection.
func NewAppSectionAdminAPI(ctx context.Context, opts app.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionAdminAPICaption,
			BasePath:     sectionAdminAPIBasePath,
			Privilege:    opts.Cfg.AppSections.AdminAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.AdminAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.AdminAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}

// NewAppSectionProvidersAPI - создаёт объект mrperms.AppSection.
func NewAppSectionProvidersAPI(ctx context.Context, opts app.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionProvidersAPICaption,
			BasePath:     sectionProvidersAPIBasePath,
			Privilege:    opts.Cfg.AppSections.ProvidersAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.ProvidersAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.ProvidersAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}

// NewAppSectionPublicAPI - создаёт объект mrperms.AppSection.
func NewAppSectionPublicAPI(ctx context.Context, opts app.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionPublicAPICaption,
			BasePath:     sectionPublicAPIBasePath,
			Privilege:    opts.Cfg.AppSections.PublicAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.PublicAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.PublicAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}

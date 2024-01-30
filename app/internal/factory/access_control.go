package factory

import (
	"context"
	"print-shop-back/config"
	"strings"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewAccessControl(ctx context.Context, cfg config.Config) (mrperms.AccessControl, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init roles and permissions for access control")

	m, err := mrperms.NewAccessControl(
		mrperms.AccessControlOptions{
			RolesDirPath:  cfg.AccessControl.Roles.DirPath,
			RolesFileType: cfg.AccessControl.Roles.FileType,
			Roles:         cfg.AccessControl.Roles.List,
			Privileges:    cfg.AccessControl.Privileges,
			Permissions:   cfg.AccessControl.Permissions,
		},
	)

	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("Registered roles: %s", strings.Join(m.RegisteredRoles(), ", "))
	logger.Info().Msgf("Guest role: %s", m.GuestRole())
	logger.Info().Msgf("Registered privileges: %s", strings.Join(m.RegisteredPrivileges(), ", "))
	logger.Debug().Msgf("Registered permissions:\n - %s", strings.Join(m.RegisteredPermissions(), ",\n - "))

	return m, nil
}

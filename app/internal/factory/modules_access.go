package factory

import (
	"print-shop-back/config"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewModulesAccess(cfg *config.Config, logger mrcore.Logger) (*mrperms.ModulesAccess, error) {
	logger.Info("Create and init roles and permissions for modules access")

	m, err := mrperms.NewModulesAccess(
		mrperms.ModulesAccessOptions{
			RolesDirPath:  cfg.ModulesAccess.Roles.DirPath,
			RolesFileType: cfg.ModulesAccess.Roles.FileType,
			Roles:         cfg.ModulesAccess.Roles.List,
			Privileges:    cfg.ModulesAccess.Privileges,
			Permissions:   cfg.ModulesAccess.Permissions,
		},
	)

	if err != nil {
		return nil, err
	}

	logger.Info("Registered roles: %s", strings.Join(m.RegisteredRoles(), ", "))
	logger.Info("Guest role: %s", m.GuestRole())
	logger.Info("Registered privileges: %s", strings.Join(m.RegisteredPrivileges(), ", "))
	logger.Debug("Registered permissions:\n - %s", strings.Join(m.RegisteredPermissions(), ",\n - "))

	return m, nil
}

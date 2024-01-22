package factory

import (
	"print-shop-back/config"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewAccessControl(cfg *config.Config, logger mrcore.Logger) (*mrperms.AccessControl, error) {
	logger.Info("Create and init roles and permissions for access control")

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

	logger.Info("Registered roles: %s", strings.Join(m.RegisteredRoles(), ", "))
	logger.Info("Guest role: %s", m.GuestRole())
	logger.Info("Registered privileges: %s", strings.Join(m.RegisteredPrivileges(), ", "))
	logger.Debug("Registered permissions:\n - %s", strings.Join(m.RegisteredPermissions(), ",\n - "))

	return m, nil
}

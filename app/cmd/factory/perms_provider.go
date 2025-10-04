package factory

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/groups"
	"github.com/mondegor/go-webcore/mraccess/role/filestorage"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

// InitPermsProvider - создаёт объект filestorage.PermsProvider.
func InitPermsProvider(logger mrlog.Logger, cfg config.Config) (*filestorage.PermsProvider, error) {
	mrlog.Info(logger, "Create and init roles and permissions for access control")

	provider, err := filestorage.NewPermsProvider(
		cfg.AccessControl.RolesDirPath,
		cfg.AccessControl.Roles,
		filestorage.WithPrivileges(cfg.AccessControl.Privileges),
		filestorage.WithPermissions(cfg.AccessControl.Permissions),
	)
	if err != nil {
		return nil, err
	}

	info := filestorage.NewRegisteredPermsInfo(provider)

	mrlog.Info(logger, "Registered roles: "+strings.Join(info.Roles, ", "))
	mrlog.Info(logger, "Registered privileges: "+strings.Join(info.Privileges, ", "))
	mrlog.Debug(logger, "Registered permissions:\n - "+strings.Join(info.Permissions, ",\n - "))

	return provider, nil
}

// InitRealmKindRights - создаёт объект groups.Groups.
func InitRealmKindRights(logger mrlog.Logger, realms []auth.UserRealm, rights mraccess.RightsSource) *groups.Groups {
	mrlog.Info(logger, "Create and init roles and permissions for access control")

	nKinds := 0
	for _, realm := range realms {
		nKinds += len(realm.UserKinds)
	}

	realmKinds := make([]groups.Group, 0, nKinds)

	for _, realm := range realms {
		for _, kind := range realm.UserKinds {
			if len(kind.Roles) == 0 {
				continue
			}

			realmKinds = append(
				realmKinds,
				groups.Group{
					Name:  realm.Name + "/" + kind.Name, // realm/kind
					Roles: kind.Roles,
				},
			)
		}
	}

	return groups.NewGroups(realmKinds, rights)
}

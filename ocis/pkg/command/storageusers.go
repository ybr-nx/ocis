// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
)

// StorageUsersCommand is the entrypoint for the storage-users command.
func StorageUsersCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-users",
		Usage:    "Start storage and data provider for /users mount",
		Category: "Extensions",
		Flags:    flagset.StorageUsersWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			scfg := configureStorageUsers(cfg)

			return cli.HandleAction(
				command.StorageUsers(scfg).Action,
				c,
			)
		},
	}
}

func configureStorageUsers(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
		cfg.Storage.Tracing.Service = cfg.Tracing.Service
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(StorageUsersCommand)
}

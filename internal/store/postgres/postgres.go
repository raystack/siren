package postgres

import (
	"github.com/goto/salt/db"
	"github.com/goto/siren/internal/store/postgres/migrations"
)

func Migrate(cfg db.Config) error {
	if err := db.RunMigrations(cfg, migrations.FS, migrations.ResourcePath); err != nil {
		return err
	}
	return nil
}

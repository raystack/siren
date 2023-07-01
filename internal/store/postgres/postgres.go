package postgres

import (
	"github.com/raystack/salt/db"
	"github.com/raystack/siren/internal/store/postgres/migrations"
)

func Migrate(cfg db.Config) error {
	if err := db.RunMigrations(cfg, migrations.FS, migrations.ResourcePath); err != nil {
		return err
	}
	return nil
}

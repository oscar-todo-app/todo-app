package mg

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	migrations "github.com/oscarsjlh/todo/migrate"
	"github.com/pressly/goose/v3"
)

func MigrateDb(dbAddr string) error {
	migration, err := goose.OpenDBWithDriver("pgx", dbAddr)
	if err != nil {
		log.Fatal("Failed migration")
		log.Fatalf(err.Error())
	}

	goose.SetBaseFS(migrations.Migrations)
	err = goose.Up(migration, ".")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return nil
}

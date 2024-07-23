package db

import (
	"database/sql"
	"embed"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

type DBClient struct {
	db     *sql.DB
	logger logrus.FieldLogger
}

func NewDBClient(sqlitePath string, logger logrus.FieldLogger) (*DBClient, error) {
	loggerWithFields := logger.WithFields(logrus.Fields{
		"package": "db",
		"struct":  "DBClient",
	})

	db, err := sql.Open("sqlite", "file:"+sqlitePath)
	if err != nil {
		log.Fatalln("failed to open sqlite connection")
	}
	if err := ensureSchema(db); err != nil {
		log.Fatalln("migration failed")
	}

	return &DBClient{
		db:     db,
		logger: loggerWithFields,
	}, nil
}

const schemaVersion = 1

func ensureSchema(db *sql.DB) error {
	sourceInstance, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return errors.Wrap(err, "invalid source instance")
	}
	targetInstance, err := sqlite.WithInstance(db, new(sqlite.Config))
	if err != nil {
		return errors.Wrap(err, "invalid source instance")
	}
	m, err := migrate.NewWithInstance(
		"httpfs", sourceInstance, "sqlite", targetInstance)
	if err != nil {
		return errors.Wrap(err, "failed to initialize migrate instance")
	}
	err = m.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}

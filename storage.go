package draken

import (
	"context"
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Storage interface {
	Init(bool)
	Stop()
	Bun() *bun.DB
	Ctx() context.Context
}

type SqlDatabase struct {
	DB      *sql.DB
	Client  *bun.DB
	Context context.Context
	Cancel  context.CancelFunc
}

func (d *Draken) initStorage() {
	if !d.Config.Storage.Enabled {
		log.Debug().Msgf("Storage is disabled in the config, skipping...")
		return
	}
	log.Debug().Msgf("Initializing storage...")

	switch d.Config.Storage.Type {
	case StorageTypeSqlite:
		d.Storage = NewSqlite("draken", "data")
	case StorageTypeLibsql:
		d.Storage = NewLibsql(d.Config.Storage.DSN)
	case StorageTypePostgres:
		d.Storage = NewPostgres(d.Config.Storage.DSN)
	}
	log.Info().Msgf("Storage initialized.")

}

func NewSqlite(sqliteFolder ...string) *SqlDatabase {
ConnectionStart:
	log.Debug().Msgf("Initializing the sqlite database...")

	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user home directory")
	}
	dataFolder := home

	// Set up sqlite connection
	for _, f := range sqliteFolder {
		dataFolder = path.Join(dataFolder, f)
	}
	err = os.MkdirAll(dataFolder, 0755)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create data directory")
	}

	conn, err := sql.Open("sqlite3", filepath.Join(dataFolder, "main.db"))
	if err != nil {
		log.Error().Msgf("Failed to connect to the database: %v", err)
		log.Warn().Msgf("Waiting for 10 seconds before trying to establish a new connection...")
		time.Sleep(10 * time.Second)
		goto ConnectionStart
	}
	bun := bun.NewDB(conn, sqlitedialect.New())

	// Create context for the database
	ctx, cancel := context.WithCancel(context.Background())
	d := &SqlDatabase{
		DB:      conn,
		Client:  bun,
		Context: ctx,
		Cancel:  cancel,
	}

	log.Info().Msgf("Initialized sqlite database.")
	return d
}

// NewPostgres creates a new postgres database
func NewPostgres(dsn string) *SqlDatabase {
ConnectionStart:
	log.Debug().Msgf("Initializing the postgres database...")

	conn := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err := conn.Ping()
	if err != nil {
		log.Error().Msgf("Failed to connect to the database: %v", err)
		log.Warn().Msgf("Waiting for 10 seconds before trying to establish a new connection...")
		time.Sleep(10 * time.Second)
		goto ConnectionStart
	}
	bun := bun.NewDB(conn, pgdialect.New())

	// Create context for the database
	ctx, cancel := context.WithCancel(context.Background())
	d := &SqlDatabase{
		DB:      conn,
		Client:  bun,
		Context: ctx,
		Cancel:  cancel,
	}

	log.Info().Msgf("Initialized postgres database.")
	return d
}

// NewLibsql creates a new libsql database
func NewLibsql(dsn string) *SqlDatabase {
ConnectionStart:
	log.Debug().Msgf("Initializing the libsql database...")

	conn, err := sql.Open("libsql", dsn)
	if err != nil {
		log.Error().Msgf("Failed to connect to the database: %v", err)
		log.Warn().Msgf("Waiting for 10 seconds before trying to establish a new connection...")
		time.Sleep(10 * time.Second)
		goto ConnectionStart
	}
	bun := bun.NewDB(conn, sqlitedialect.New())

	// Create context for the database
	ctx, cancel := context.WithCancel(context.Background())
	d := &SqlDatabase{
		DB:      conn,
		Client:  bun,
		Context: ctx,
		Cancel:  cancel,
	}

	log.Info().Msgf("Initialized libsql database.")
	return d
}

func (d *SqlDatabase) Init(debug bool) {
	d.Client.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(debug),
		bundebug.WithEnabled(debug),
		bundebug.WithWriter(log.Logger),
	))
}

func (d *SqlDatabase) Stop() {
	d.Client.Close()
}

func (d *SqlDatabase) Bun() *bun.DB {
	return d.Client
}

func (d *SqlDatabase) Ctx() context.Context {
	return d.Context
}

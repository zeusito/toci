package db

import (
	"fmt"
	"time"

	"github.com/zeusito/toci/pkg/config"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun/extra/bunzerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type DatabaseConnection struct {
	Conn *bun.DB
}

func MustCreatePooledConnection(dbConfig config.DatabaseConfigurations) *DatabaseConnection {
	if !dbConfig.Enabled {
		log.Warn().Msg("database is disabled")
		return &DatabaseConnection{}
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)

	parsedCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatal().Msgf("Error parsing database configuration: %v", err)
		return nil
	}

	// Init a connection compatible with the standard library
	dbPool := stdlib.OpenDB(*parsedCfg)

	// Connection pool settings
	dbPool.SetMaxOpenConns(dbConfig.PoolSize)
	dbPool.SetMaxIdleConns(dbConfig.PoolSize)
	dbPool.SetConnMaxLifetime(3 * time.Minute)

	// Test connection
	if err := dbPool.Ping(); err != nil {
		log.Fatal().Msgf("Error pinging database: %v", err)
		return nil
	}

	log.Info().Msg("Successfully connected to database")

	hook := bunzerolog.NewQueryHook(
		bunzerolog.WithQueryLogLevel(zerolog.DebugLevel),
		bunzerolog.WithSlowQueryLogLevel(zerolog.WarnLevel),
		bunzerolog.WithErrorQueryLogLevel(zerolog.ErrorLevel),
		bunzerolog.WithSlowQueryThreshold(3*time.Second),
	)

	db := bun.NewDB(dbPool, pgdialect.New(), bun.WithDiscardUnknownColumns()).
		WithQueryHook(hook)

	return &DatabaseConnection{Conn: db}
}

func (c *DatabaseConnection) Close() {
	if c.Conn == nil {
		return
	}
	_ = c.Conn.Close()
}

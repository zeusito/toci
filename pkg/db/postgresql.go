package db

import (
	"context"
	"fmt"
	"time"

	"github.com/zeusito/toci/pkg/config"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun/extra/bunzerolog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type DatabaseConnection struct {
	Conn *bun.DB
	pool *pgxpool.Pool
}

func MustCreatePooledConnection(dbConfig config.DatabaseConfigurations) *DatabaseConnection {
	if !dbConfig.Enabled {
		log.Warn().Msg("database is disabled")
		return &DatabaseConnection{}
	}

	// Format the DSN, it aims for a fixed number of connections
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s&pool_max_conns=%d&pool_min_conns=%d",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName,
		dbConfig.PoolSize, dbConfig.PoolSize)

	// Parse the DSN
	parsedCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing database configuration")
		return nil
	}

	// Create the connection pool using pgxpool
	pool, err := pgxpool.NewWithConfig(context.Background(), parsedCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating connection pool")
		return nil
	}

	// Init a connection compatible with the standard library
	dbPool := stdlib.OpenDBFromPool(pool)

	// Test connection
	if err := dbPool.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Error pinging database")
		return nil
	}

	log.Info().Msgf("Successfully connected to database. Pool size: %d", pool.Stat().TotalConns())

	hook := bunzerolog.NewQueryHook(
		bunzerolog.WithQueryLogLevel(zerolog.InfoLevel),
		bunzerolog.WithSlowQueryLogLevel(zerolog.WarnLevel),
		bunzerolog.WithErrorQueryLogLevel(zerolog.ErrorLevel),
		bunzerolog.WithSlowQueryThreshold(3*time.Second),
	)

	db := bun.NewDB(dbPool, pgdialect.New(), bun.WithDiscardUnknownColumns()).
		WithQueryHook(hook)

	return &DatabaseConnection{Conn: db, pool: pool}
}

func (c *DatabaseConnection) Close() {
	if c.pool == nil {
		return
	}

	c.pool.Close()
}

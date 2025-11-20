package database

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/WahyuSiddarta/be_saham_go/config"

	// PostgreSQL DB Driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PSQLGetDBReadWrite : init DB with optimization layer
func PSQLGetDBReadWrite() *sqlx.DB {

	Logger.Info().Msg("PQRW - DB open started")
	cfg := config.Get()
	if cfg == nil {
		Logger.Error().Msg("PQRW - Configuration is nil")
		return nil
	}

	// Get database configuration for read-write connection
	dbConfig := cfg.Database.RW

	// Debug logging to see what values we got
	Logger.Info().Str("host", dbConfig.Host).Str("port", dbConfig.Port).Str("user", dbConfig.User).Str("db", dbConfig.DBName).Str("schema", dbConfig.Schema).Msg("PQRW - DB connection parameters")

	if dbConfig.Host == "" {
		Logger.Fatal().Msg("PQRW - DB host is empty, check configuration")
		return nil
	}
	if dbConfig.Port == "" {
		Logger.Fatal().Msg("PQRW - DB port is empty, check configuration")
		return nil
	}

	db, err := sqlx.Open("pgx",
		"host="+dbConfig.Host+" port="+dbConfig.Port+" user="+dbConfig.User+" password="+dbConfig.Password+" dbname="+dbConfig.DBName+"")

	if err != nil {
		Logger.Fatal().Err(err).Str("host", dbConfig.Host).Msg("PQRW - DB open failed")
		return nil
	}

	err = db.Ping()
	if err != nil {
		Logger.Fatal().Err(err).Str("host", dbConfig.Host).Msg("PQRW - DB ping failed")
		return nil
	}

	// Set search_path after connection if needed
	if dbConfig.Schema != "" {
		_, err = db.Exec("SET search_path TO " + dbConfig.Schema)
		if err != nil {
			Logger.Fatal().Err(err).Str("schema", dbConfig.Schema).Msg("PQRW - Failed to set search_path")
			return nil
		}
	}

	// Enhanced connection pool settings for write operations
	if dbConfig.MaxCon > 0 {
		db.SetMaxOpenConns(dbConfig.MaxCon)
	} else {
		db.SetMaxOpenConns(25) // Default for read-write operations
	}

	if dbConfig.MaxIdle > 0 {
		db.SetMaxIdleConns(dbConfig.MaxIdle)
	} else {
		db.SetMaxIdleConns(10) // Default for better connection reuse
	}

	// Optimized connection lifetime and idle timeout
	db.SetConnMaxLifetime(45 * time.Minute) // Reduced from 1 hour for better connection cycling
	db.SetConnMaxIdleTime(10 * time.Minute) // New: Close idle connections after 10 minutes

	Logger.Info().Msg("PQRW - DB open completed with optimization layer")
	return db
}

// PSQLGetDBReadCache : init read-cache DB connection with optimization
func PSQLGetDBReadCache() *sqlx.DB {

	Logger.Info().Msg("PQRC - DB open started")
	cfg := config.Get()
	if cfg == nil {
		Logger.Error().Msg("PQRC - Configuration is nil")
		return nil
	}

	// Get database configuration for read-cache connection
	dbConfig := cfg.Database.RC

	// Debug logging to see what values we got
	Logger.Info().Str("host", dbConfig.Host).Str("port", dbConfig.Port).Str("user", dbConfig.User).Str("db", dbConfig.DBName).Str("schema", dbConfig.Schema).Msg("PQRC - DB connection parameters")

	if dbConfig.Host == "" {
		Logger.Fatal().Msg("PQRC - DB host is empty, check configuration")
		return nil
	}
	if dbConfig.Port == "" {
		Logger.Fatal().Msg("PQRC - DB port is empty, check configuration")
		return nil
	}

	db, err := sqlx.Open("pgx",
		"host="+dbConfig.Host+" port="+dbConfig.Port+" user="+dbConfig.User+" password="+dbConfig.Password+" dbname="+dbConfig.DBName+"")
	if err != nil {
		Logger.Fatal().Err(err).Str("host", dbConfig.Host).Msg("PQRC - DB open failed")
		return nil
	}

	err = db.Ping()
	if err != nil {
		Logger.Fatal().Err(err).Str("host", dbConfig.Host).Msg("PQRC - DB ping failed")
		return nil
	}

	// Set search_path after connection if needed
	if dbConfig.Schema != "" {
		_, err = db.Exec("SET search_path TO " + dbConfig.Schema)
		if err != nil {
			Logger.Fatal().Err(err).Str("schema", dbConfig.Schema).Msg("PQRC - Failed to set search_path")
			return nil
		}
	}

	// Enhanced connection pool settings for cache read operations
	if dbConfig.MaxCon > 0 {
		db.SetMaxOpenConns(dbConfig.MaxCon)
	} else {
		db.SetMaxOpenConns(25) // Default for cache operations
	}

	if dbConfig.MaxIdle > 0 {
		db.SetMaxIdleConns(dbConfig.MaxIdle)
	} else {
		db.SetMaxIdleConns(10) // Default idle connections for cache
	}

	Logger.Info().Msg("PQRC - DB open completed (using pgbouncer)")
	return db
}

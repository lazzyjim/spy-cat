package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"spy-cat/src/config"
	"time"
)

var Connection *sql.DB

func Init(cof *config.Postgres) error {
	conn, err := initializeDatabaseSession(cof)
	if err != nil {
		return fmt.Errorf("can not connect to DB DB, connection string %s err: %s", cof.ConnectionSource(), err.Error())
	}
	Connection = conn
	return nil
}

func initializeDatabaseSession(cof *config.Postgres) (*sql.DB, error) {
	db, err := sql.Open("postgres", cof.ConnectionSource())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can not ping DB connection, source string %s, err: %s", cof.ConnectionSource(), err)
	}
	// @link https://www.alexedwards.net/blog/configuring-sqldb
	// Set the maximum number of concurrently open connections to 5. Setting this
	// to less than or equal to 0 will mean there is no maximum limit (which
	// is also the default setting).
	db.SetMaxOpenConns(cof.Settings.MaximumPoolSize)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(cof.Settings.MaximumPoolSize)
	// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(time.Duration(cof.Settings.ConnectionTimeout) * time.Second)

	return db, nil
}

func Conn() *sql.DB {
	return Connection
}

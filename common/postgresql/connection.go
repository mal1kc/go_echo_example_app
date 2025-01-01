package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

func GetConnectionPool(context context.Context, connectionString string) *pgxpool.Pool {
	connConfig, parseConfigErr := pgxpool.ParseConfig(connectionString)
	if parseConfigErr != nil {
		panic(parseConfigErr)
	}

	conn, err := pgxpool.ConnectConfig(context, connConfig)
	if err != nil {
		log.Errorf("unable to connect to database: %s\n", err)
		panic(err)
	}
	return conn
}

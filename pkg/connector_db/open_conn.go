package connector_db

import (
	"auth/pkg/connector_db/interface"
	"auth/pkg/connector_db/postgre"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func openPostger(config string) _interface.DB {
	connConfig, _ := pgxpool.ParseConfig(config)
	conn, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("error:", err)
	}
	return postgre.NewPostreConn(conn)
}

func Open(typedb, config string) _interface.DB {
	switch typedb {
	case "postgre":
		return openPostger(config)
	default:
		return nil
	}
}

package repo

import (
	"context"
	"fmt"
	"magnusquiz/pkg/log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var connPool *pgxpool.Pool

func Connect(host, user, pass, name string, port uint16) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		user,
		pass,
		host,
		port,
		name,
	)
	for {
		pool, err := pgxpool.Connect(context.Background(), dsn)
		if err != nil {
			log.Logger.Error(err)
			time.Sleep(time.Second)
			continue
		}

		if err := pool.Ping(context.Background()); err != nil {
			log.Logger.Error(err)
			time.Sleep(time.Second)
			continue
		}
		connPool = pool
		break

	}
}

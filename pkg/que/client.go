package que

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/vgarvardt/gue/v2"
	"github.com/vgarvardt/gue/v2/adapter/pgxv4"
)

func getClient(mypool *pgxpool.Pool) (*gue.Client, error) {
	poolAdapter := pgxv4.NewConnPool(mypool)
	gc := gue.NewClient(poolAdapter)
	return gc, nil
}

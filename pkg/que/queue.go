package que

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/pkg/config"
	"github.com/vgarvardt/gue/v2"
)

type Que struct {
	client  *gue.Client
	worker  *gue.WorkerPool
	quename string
}

func Get(db *pgxpool.Pool, cfg *config.QueConfig) (*Que, error) {

	clnt, err := getClient(db)
	if err != nil {
		return nil, err
	}

	var wrk *gue.WorkerPool

	if cfg.WorkerEnabled {
		wrk, err = getWorkerpool(clnt, cfg)
		if err != nil {
			return nil, err
		}

	} else {
		wrk = nil
	}

	return &Que{
		client:  clnt,
		worker:  wrk,
		quename: cfg.QueName,
	}, nil
}

func (c *Que) Enquejob(j *gue.Job) error {
	j.Queue = c.quename
	return c.client.Enqueue(context.Background(), j)
}

func (c *Que) StartWorkerPool() error {
	fmt.Println(c.worker != nil)
	if c.worker != nil {
		fmt.Println("worker started")
		return c.worker.Run(context.Background())
	}
	return errors.New("Worker not allowed")
}

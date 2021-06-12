package que

import (
	"time"

	"github.com/sveltegobackend/pkg/config"
	"github.com/vgarvardt/gue/v2"
)

func getWorkerpool(gc *gue.Client, cfg *config.QueConfig) (*gue.WorkerPool, error) {
	duration, _ := time.ParseDuration("15s")
	workers := gue.NewWorkerPool(gc, wm, int(cfg.WorkerCount), gue.WithPoolQueue(cfg.QueName), gue.WithPoolPollInterval(duration))
	return workers, nil
}

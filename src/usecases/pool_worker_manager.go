package usecases

import (
	"os"
	"runtime"
	"strconv"

	"github.com/robertokbr/blinkchat/src/domain/models"
)

type PoolWorkerManager struct {
	pool *models.Pool
	jobs chan models.Message
}

func NewPoolWorkerManager(pool *models.Pool, jobs chan models.Message) *PoolWorkerManager {
	return &PoolWorkerManager{
		pool: pool,
		jobs: jobs,
	}
}

func (uc *PoolWorkerManager) Stop() {
	close(uc.jobs)
	close(uc.pool.Pairs)
}

func (uc *PoolWorkerManager) Start() {
	numberOfWorkers, err := strconv.Atoi(os.Getenv("NUMBER_OF_WORKERS"))

	if err != nil {
		numberOfWorkers = runtime.NumCPU()
	}

	simpleMatchPoolPairs := NewSimpleMatchPoolPairs(uc.pool)

	for i := 0; i < numberOfWorkers; i++ {
		go PoolWorker(i, uc.pool, uc.jobs, simpleMatchPoolPairs)
	}
}

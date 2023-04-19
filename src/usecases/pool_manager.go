package usecases

type PoolManager struct {
	pooWorker      *PoolWorker
	matchPoolPairs *MatchPoolPairs
}

func NewPoolManager(poolWorker *PoolWorker, matchPoolPairs *MatchPoolPairs) *PoolManager {
	return &PoolManager{
		pooWorker:      poolWorker,
		matchPoolPairs: matchPoolPairs,
	}
}

func (uc *PoolManager) Execute(numberOfWorkers int) {
	go uc.matchPoolPairs.Execute()

	for i := 0; i < numberOfWorkers; i++ {
		go uc.pooWorker.Execute(i)
	}
}

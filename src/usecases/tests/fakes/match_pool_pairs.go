package usecases_tests_fakes

type MatchPoolPairs struct{}

func NewMatchPoolPairs() *MatchPoolPairs {
	return &MatchPoolPairs{}
}

func (mpp *MatchPoolPairs) Execute() {}

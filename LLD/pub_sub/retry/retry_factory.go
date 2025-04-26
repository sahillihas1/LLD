package retry

import "github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"

type IFactory interface {
	GetRetryStrategy(strategy model.RetryStrategy) IRetry
}

type Factory struct {
}

func NewRetryFactory() *Factory {
	return &Factory{}
}

func (r *Factory) GetRetryStrategy(strategy model.RetryStrategy) IRetry {
	switch strategy {
	case model.Linear:
		return NewLinearRetry()
	case model.Exponential:
		return NewExponentialRetry()
	}
	return nil
}

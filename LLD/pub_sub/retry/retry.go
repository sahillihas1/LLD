package retry

import (
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"
	"time"
)

type IRetry interface {
	RetryMessage(subscriber *model.Subscriber, message *model.Message) error
}

type LinearRetry struct {
}

func NewLinearRetry() *LinearRetry {
	return &LinearRetry{}
}

func (l *LinearRetry) RetryMessage(subscriber *model.Subscriber, message *model.Message) error {
	if message.RetryCount < subscriber.TotalRetry {
		message.RetryCount++
		time.AfterFunc(time.Second*time.Duration(message.RetryCount), func() {
			sendMessageToSub(subscriber, message)
		})
	}
	return nil
}

func sendMessageToSub(subscriber *model.Subscriber, message *model.Message) error {
	subscriber.Ch <- *message
	return nil
}

type ExponentialRetry struct {
}

func NewExponentialRetry() *ExponentialRetry {
	return &ExponentialRetry{}
}

func (e *ExponentialRetry) RetryMessage(subscriber *model.Subscriber, message *model.Message) error {
	return nil
}

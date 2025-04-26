package services

import (
	"fmt"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/retry"
)

type ISubscriberService interface {
	CreateSubscriber(id int) *model.Subscriber
	ConsumerMessage(s *model.Subscriber)
}

type SubscriberService struct {
	retryFactory retry.IFactory
}

func NewSubscriberService() *SubscriberService {
	return &SubscriberService{
		retryFactory: retry.NewRetryFactory(),
	}
}

func (s SubscriberService) CreateSubscriber(id int) *model.Subscriber {
	return &model.Subscriber{
		ID:         id,
		Ch:         make(chan model.Message, 100), // Buffered channel for messages
		Done:       make(chan bool),
		TotalRetry: 3,
	}
}

func (s SubscriberService) ConsumerMessage(a *model.Subscriber) {
	go func() {
		for {
			select {
			case msg := <-a.Ch:
				if msg.Content == "Hello World" {
					s.retryFactory.GetRetryStrategy(a.RetryStrategy).RetryMessage(a, &msg)
				}
				fmt.Printf("Subscriber %d received message: %s\n", a.ID, msg.Content)
			case <-a.Done:
				fmt.Printf("Subscriber %d stopping.\n", a.ID)
				return
			}
		}
	}()
}

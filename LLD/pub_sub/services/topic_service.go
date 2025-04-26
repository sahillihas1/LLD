package services

import (
	"fmt"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/repository"
	"sync"
)

type ITopicService interface {
	CreateTopic(topic *model.Topic) error
	AddSubscriber(topic string, subscriber *model.Subscriber) error
	Publish(topic string, msg model.Message) error
}

type TopicService struct {
	repo repository.ITopicRepository
	mu   sync.Mutex
}

func NewTopicService(repo repository.ITopicRepository) *TopicService {
	return &TopicService{
		repo: repo,
	}
}

func (ts *TopicService) CreateTopic(topic *model.Topic) error {
	return ts.repo.CreateTopic(topic)
}

func (ts *TopicService) AddSubscriber(topic string, subscriber *model.Subscriber) error {
	return ts.repo.AddSubscriber(topic, subscriber)
}

func (t *TopicService) Publish(topic string, msg model.Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	topicDetail, err := t.repo.GetTopic(topic)
	if err != nil {
		return err
	}
	for _, subscriber := range topicDetail.Subscribers {
		subscriber.Ch <- msg
	}
	fmt.Printf("Message published to topic %s: %s\n", topicDetail.Name, msg.Content)
	return nil
}

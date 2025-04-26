package repository

import "github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"

type ITopicRepository interface {
	CreateTopic(topic *model.Topic) error
	GetTopic(topic string) (*model.Topic, error)
	AddSubscriber(topic string, subscriber *model.Subscriber) error
}

type TopicRepository struct {
	topics map[string]*model.Topic
}

func NewTopicRepository() *TopicRepository {
	return &TopicRepository{
		topics: make(map[string]*model.Topic),
	}
}

func (tr *TopicRepository) CreateTopic(topic *model.Topic) error {
	tr.topics[topic.Name] = topic
	return nil
}

func (tr *TopicRepository) GetTopic(topic string) (*model.Topic, error) {
	if topic, ok := tr.topics[topic]; ok {
		return topic, nil
	}
	return nil, nil
}

func (tr *TopicRepository) AddSubscriber(topic string, subscriber *model.Subscriber) error {
	if topic, ok := tr.topics[topic]; ok {
		topic.Subscribers[subscriber.ID] = subscriber
		return nil
	}
	return nil
}

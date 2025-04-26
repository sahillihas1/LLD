package main

import (
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/model"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/repository"
	"github.com/SahilSrivastava/Downloads/machinecoding/cache_system/services"
	"time"
)

type PubSubService struct {
	TopicSerivce      services.ITopicService
	SubscriberService services.ISubscriberService
}

func NewPubSubService(topicService services.ITopicService, subscriberService services.ISubscriberService) *PubSubService {
	return &PubSubService{
		TopicSerivce:      topicService,
		SubscriberService: subscriberService,
	}
}

func main() {
	topicRepo := repository.NewTopicRepository()
	topicService := services.NewTopicService(topicRepo)
	subscriberService := services.NewSubscriberService()
	pubSubService := NewPubSubService(topicService, subscriberService)
	pubSubService.TopicSerivce.CreateTopic(model.NewTopic("topic1"))
	sub1 := pubSubService.SubscriberService.CreateSubscriber(1)
	pubSubService.TopicSerivce.AddSubscriber("topic1", sub1)
	pubSubService.TopicSerivce.Publish("topic1", model.Message{Content: "Hello World"})
	pubSubService.TopicSerivce.Publish("topic1", model.Message{Content: "Hello World2"})
	pubSubService.SubscriberService.ConsumerMessage(sub1)
	pubSubService.TopicSerivce.Publish("topic1", model.Message{Content: "Hello World3"})
	time.Sleep(10 * time.Second)
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// --- Models ---

type Message struct {
	Offset  int
	Content string
}

type Subscriber struct {
	ID            int
	CurrentOffset int
	Done          chan struct{}
	offsetLock    sync.Mutex
}

func (s *Subscriber) SetOffset(offset int) {
	s.offsetLock.Lock()
	defer s.offsetLock.Unlock()
	s.CurrentOffset = offset
	fmt.Printf("Subscriber %d manually set offset to %d\n", s.ID, offset)
}

type Topic struct {
	Name        string
	Messages    []Message
	Subscribers []*Subscriber
}

// --- Interfaces ---

type ITopicService interface {
	CreateTopic(topic *Topic) error
	AddSubscriber(topic string, subscriber *Subscriber) error
	RemoveSubscriber(topicName string, subscriber *Subscriber) error
	Publish(topic string, content string) error
}

type ISubscriberService interface {
	CreateSubscriber(id int) *Subscriber
	ConsumeMessages(s *Subscriber, topic *Topic)
}

// --- Services ---

type TopicService struct {
	topics            map[string]*Topic
	topicLock         sync.RWMutex
	subscriberService ISubscriberService
}

func NewTopicService(subService ISubscriberService) ITopicService {
	return &TopicService{
		topics:            make(map[string]*Topic),
		subscriberService: subService,
	}
}

func (ts *TopicService) CreateTopic(topic *Topic) error {
	ts.topicLock.Lock()
	defer ts.topicLock.Unlock()

	if _, exists := ts.topics[topic.Name]; exists {
		return fmt.Errorf("topic already exists")
	}
	ts.topics[topic.Name] = topic
	return nil
}

func (ts *TopicService) AddSubscriber(topicName string, subscriber *Subscriber) error {
	ts.topicLock.Lock()
	defer ts.topicLock.Unlock()

	topic, exists := ts.topics[topicName]
	if !exists {
		return fmt.Errorf("topic not found")
	}
	topic.Subscribers = append(topic.Subscribers, subscriber)

	go ts.subscriberService.ConsumeMessages(subscriber, topic)

	return nil
}

func (ts *TopicService) RemoveSubscriber(topicName string, subscriber *Subscriber) error {
	ts.topicLock.Lock()
	defer ts.topicLock.Unlock()

	topic, exists := ts.topics[topicName]
	if !exists {
		return fmt.Errorf("topic not found")
	}
	for i, sub := range topic.Subscribers {
		if sub.ID == subscriber.ID {
			topic.Subscribers = append(topic.Subscribers[:i], topic.Subscribers[i+1:]...)
			select {
			case <-sub.Done:
			default:
				close(sub.Done)
			}

			return nil
		}
	}

	return fmt.Errorf("subscriber not found in topic")
}

func (ts *TopicService) Publish(topicName string, content string) error {
	ts.topicLock.Lock()
	defer ts.topicLock.Unlock()

	topic, exists := ts.topics[topicName]
	if !exists {
		return fmt.Errorf("topic not found")
	}

	newMsg := Message{
		Offset:  len(topic.Messages),
		Content: content,
	}
	topic.Messages = append(topic.Messages, newMsg)

	return nil
}

type SubscriberService struct{}

func NewSubscriberService() ISubscriberService {
	return &SubscriberService{}
}

func (ss *SubscriberService) CreateSubscriber(id int) *Subscriber {
	return &Subscriber{
		ID:            id,
		CurrentOffset: 0,
		Done:          make(chan struct{}),
	}
}

func (ss *SubscriberService) ConsumeMessages(s *Subscriber, topic *Topic) {
	for {
		select {
		case <-s.Done:
			fmt.Printf("Subscriber %d stopping consumption.\n", s.ID)
			return
		default:
			s.offsetLock.Lock()
			if s.CurrentOffset < len(topic.Messages) {
				msg := topic.Messages[s.CurrentOffset]
				//	s.Channel <- msg
				fmt.Printf("Subscriber %d received [offset %d]: %s\n", s.ID, msg.Offset, msg.Content)
				s.CurrentOffset++
				s.offsetLock.Unlock()
			} else {
				s.offsetLock.Unlock()
				time.Sleep(500 * time.Millisecond) // Wait for new messages
			}
		}
	}
}

// --- Main Demo ---

func main() {
	// Setup Services
	subscriberService := NewSubscriberService()
	topicService := NewTopicService(subscriberService)

	// Create Topic
	topic := &Topic{Name: "technology"}
	_ = topicService.CreateTopic(topic)

	// Create Subscriber
	sub := subscriberService.CreateSubscriber(1)
	subb := subscriberService.CreateSubscriber(2)

	// Add Subscriber to Topic
	_ = topicService.AddSubscriber("technology", sub)
	_ = topicService.AddSubscriber("technology", subb)

	// Publish Messages
	_ = topicService.Publish("technology", "Message 1: New AI model released!")
	time.Sleep(500 * time.Millisecond)
	_ = topicService.Publish("technology", "Message 2: Quantum computing breakthrough!")
	time.Sleep(1000 * time.Millisecond)
	_ = topicService.RemoveSubscriber("technology", subb)
	_ = topicService.Publish("technology", "Message 3: Self-driving cars 2.0 announced!")
	time.Sleep(2 * time.Second)

	//---- Change Offset Manually ----
	fmt.Println("=== Resetting offset to 1 ===")
	sub.SetOffset(1)

	time.Sleep(5 * time.Second) // Let consumer reconsume from offset 1
}

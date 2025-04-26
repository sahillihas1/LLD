package model

import "sync"

type Topic struct {
	Name        string
	Subscribers map[int]*Subscriber
	Mu          sync.Mutex
	Messages    []*Message
	Cond        *sync.Cond
}

func NewTopic(name string) *Topic {
	mu := &sync.Mutex{}
	return &Topic{
		Name:        name,
		Messages:    make([]*Message, 0), // Initialize message store
		Subscribers: make(map[int]*Subscriber),
		Mu:          *mu,
		Cond:        sync.NewCond(mu), // Create the condition variable for notifying
	}
}

type Subscriber struct {
	ID            int
	Ch            chan Message
	Done          chan bool
	consume       func(msg *Message)
	TotalRetry    int
	RetryStrategy RetryStrategy
}

type RetryStrategy int

const (
	Linear RetryStrategy = iota
	Exponential
)

type Message struct {
	Content    string
	RetryCount int
}

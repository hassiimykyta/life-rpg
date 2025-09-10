package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type ProducerFactory struct {
	brokers      []string
	mu           sync.RWMutex
	cache        map[string]*Producer
	requiredAcks kafka.RequiredAcks
	balancer     kafka.Balancer
}

type ProducerFactoryConfig struct {
	Brokers      []string
	RequiredAcks kafka.RequiredAcks
	Balancer     kafka.Balancer
}

func NewProducerFactory(cfg ProducerFactoryConfig) *ProducerFactory {
	return &ProducerFactory{
		brokers:      cfg.Brokers,
		cache:        make(map[string]*Producer),
		requiredAcks: cfg.RequiredAcks,
		balancer:     cfg.Balancer,
	}
}

func (f *ProducerFactory) Get(topic string) *Producer {
	f.mu.RLock()
	if p, ok := f.cache[topic]; ok {
		f.mu.RUnlock()
		return p
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.cache[topic]; ok {
		return p
	}

	p := NewProducer(ProducerConfig{
		Brokers:      f.brokers,
		Topic:        topic,
		RequiredAcks: f.requiredAcks,
		Balancer:     f.balancer,
	})
	f.cache[topic] = p
	return p
}

func (f *ProducerFactory) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var firstErr error
	for topic, p := range f.cache {
		if err := p.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(f.cache, topic)
	}
	return firstErr
}

package worker

import (
	"context"
	"sync"
)

type TaskProcessor interface {
	Process(ctx context.Context, payload []byte) ([]byte, error)
}

type Registry struct {
	mu         sync.RWMutex
	processors map[string]TaskProcessor
}

func NewRegistry() *Registry {
	return &Registry{
		processors: make(map[string]TaskProcessor),
	}
}

// добавляет новый TaskProcessor в реестр.
func (r *Registry) Register(taskType string, processor TaskProcessor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processors[taskType] = processor
}

// возвращает TaskProcessor по его типу.
func (r *Registry) GetProcessor(taskType string) TaskProcessor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.processors[taskType]
}
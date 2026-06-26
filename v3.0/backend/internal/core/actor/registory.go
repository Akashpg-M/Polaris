package actor

import (
	"sync"
)

// ActorRegistry maps active hardware IDs to stateful actor runtimes safely
type ActorRegistry struct {
	mu        sync.RWMutex
	actors    map[string]*AssetActor
	publisher EventPublisher
	capacity  int
}

func NewActorRegistry(pub EventPublisher, mailboxCapacity int) *ActorRegistry {
	return &ActorRegistry{
		actors:    make(map[string]*AssetActor),
		publisher: pub,
		capacity:  mailboxCapacity,
	}
}

// GetOrCreate retrieves an active actor runline or instantiates a fresh context boundary
func (r *ActorRegistry) GetOrCreate(id string) *AssetActor {
	r.mu.RLock()
	act, exists := r.actors[id]
	r.mu.RUnlock()

	if exists {
		return act
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check check boundary under write lock safety
	if act, exists = r.actors[id]; exists {
		return act
	}

	newActor := NewAssetActor(id, r.publisher, r.capacity)
	newActor.Start()
	r.actors[id] = newActor
	return newActor
}

// Remove shuts down an actor lifecycle context cleanly
func (r *ActorRegistry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if act, exists := r.actors[id]; exists {
		act.Stop()
		delete(r.actors, id)
	}
}
package broadcast

import "sync"

// Synced broadcaster provide broadcast without running a goroutine in the background.
// It should act the same as UnblockedbroadCaster.
// The zero SyncedBroadcaster is empty and ready for use. A broadcaster must not be copied after first use.
type SyncedBroadcaster[T any] struct {
	m sync.Map
}

func (b *SyncedBroadcaster[T]) Register(ch chan<- T) {
	b.m.Store(ch, struct{}{})
}

func (b *SyncedBroadcaster[T]) Unregister(ch chan<- T) {
	b.m.Delete(ch)
}

func (b *SyncedBroadcaster[T]) Subbmit(m T) {
	b.m.Range(func(key, _ any) bool {
		select {
		case key.(chan<- T) <- m:
		default:
		}
		return true
	})
}

func (b *SyncedBroadcaster[T]) Close() {
	b.m = sync.Map{}
}

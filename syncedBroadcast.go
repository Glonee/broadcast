package broadcast

import "sync"

//Synced broadcaster provide broadcast without running a goroutine in the background.
//It should act the same as UnblockedbroadCaster.
type SyncedBroadcaster[T any] struct {
	sync.Map
}

func (b *SyncedBroadcaster[T]) Register(ch chan<- T) {
	b.Store(ch, struct{}{})
}
func (b *SyncedBroadcaster[T]) Unregister(ch chan<- T) {
	b.Delete(ch)
}
func (b *SyncedBroadcaster[T]) Subbmit(m T) {
	b.Range(func(key, _ any) bool {
		select {
		case key.(chan<- T) <- m:
		default:
		}
		return true
	})
}
func (b *SyncedBroadcaster[T]) Close() {
	*b = SyncedBroadcaster[T]{}
}

package broadcast

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSyncedBroadcaster(t *testing.T) {
	b := &SyncedBroadcaster[int]{}
	defer b.Close()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int, 1)
		b.Register(cch)
		go func(c chan int) {
			defer wg.Done()
			ctx, calcel := context.WithTimeout(context.Background(), time.Second)
			defer calcel()
			select {
			case <-ctx.Done():
				t.Error("a subscriber did't receive message")
			case <-c:
			}
		}(cch)
	}
	cch := make(chan int, 1)
	b.Register(cch)
	b.Unregister(cch)
	wg.Add(1)
	go func(c chan int) {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		select {
		case <-ctx.Done():
		case <-c:
			t.Error("received message after unregister")
		}
	}(cch)
	b.Subbmit(1)
	wg.Wait()
}

func TestSyncedBroadcasterUnblock(t *testing.T) {
	b := SyncedBroadcaster[int]{}
	cch := make(chan int)
	b.Register(cch)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		b.Subbmit(1)
		b.Subbmit(1)
		cancel()
	}()
	<-ctx.Done()
	if ctx.Err() != context.Canceled {
		t.Error("broadcaster blocked")
	}
}

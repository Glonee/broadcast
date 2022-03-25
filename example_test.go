package broadcast

import (
	"sync"
	"testing"
)

func Test(t *testing.T) {
	b := NewbroadCaster[int](100)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int)
		b.Register(cch)
		go func() {
			<-cch
			b.Unregister(cch)
			wg.Done()
		}()
	}
	b.Subbmit(1)
	wg.Wait()
	//test for unblockedbroadcaster
	b = NewUnblockedbroadCaster[int](100)
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int)
		b.Register(cch)
		go func() {
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int, 1)
		b.Register(cch)
		go func() {
			<-cch
			wg.Done()
		}()
	}
	b.Subbmit(1)
	wg.Wait()
}

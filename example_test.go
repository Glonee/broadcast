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
			wg.Done()
			b.Unregister(cch)
		}()
	}
	b.Subbmit(1)
	wg.Wait()
}

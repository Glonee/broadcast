package broadcast

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	b := NewbroadCaster[int](100)
	defer b.Close()
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	finished := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			t.Error(errors.New("not all subscribers received messages"))
		case <-finished:
		}
		quit <- struct{}{}
	}()
	go func() {
		wg.Wait()
		finished <- struct{}{}
	}()
	<-quit
}
func Test_unblockedBroadcaster(t *testing.T) {
	b := NewUnblockedbroadCaster[int](100)
	defer b.Close()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		cch := make(chan int)
		b.Register(cch)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int, 1)
		b.Register(cch)
		go func() {
			<-cch
			b.Unregister(cch)
			wg.Done()
		}()
	}
	b.Subbmit(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	finished := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			t.Error(errors.New("not all subscribers received messages"))
		case <-finished:
		}
		quit <- struct{}{}
	}()
	go func() {
		wg.Wait()
		finished <- struct{}{}
	}()
	<-quit
}
func Test_SyncedBroadCaster(t *testing.T) {
	b := SyncedBroadcaster[int]{}
	defer b.Close()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		cch := make(chan int)
		b.Register(cch)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		cch := make(chan int, 1)
		b.Register(cch)
		go func() {
			<-cch
			b.Unregister(cch)
			wg.Done()
		}()
	}
	b.Subbmit(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	finished := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			t.Error(errors.New("not all subscribers received messages"))
		case <-finished:
		}
		quit <- struct{}{}
	}()
	go func() {
		wg.Wait()
		finished <- struct{}{}
	}()
	<-quit
}

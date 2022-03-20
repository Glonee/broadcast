package broadcast

type Operate[T any] struct {
	Chan chan<- T
	Done chan<- struct{}
}
type broadcaster[T any] struct {
	input     chan T
	reg       chan Operate[T]
	unreg     chan Operate[T]
	registers map[chan<- T]struct{}
}
type BroadCaster[T any] interface {
	//Register a new channel to receive messages.
	Register(chan<- T)
	//Unregister a channel that no longer receive messages.
	Unregister(chan<- T)
	//Subbmit a new message to all subscribers.
	Subbmit(T)
	//Close the broadcaster.
	Close()
}

//Create a new broadcaster with the given input channel buffer length.
func NewbroadCaster[T any](buflen int) BroadCaster[T] {
	b := &broadcaster[T]{
		input:     make(chan T, buflen),
		reg:       make(chan Operate[T]),
		unreg:     make(chan Operate[T]),
		registers: make(map[chan<- T]struct{}),
	}
	go b.run()
	return b
}
func (b *broadcaster[T]) run() {
	for {
		select {
		//Send message to all subscribers.
		case m := <-b.input:
			for ch := range b.registers {
				ch <- m
			}
		//Add a new subscriber.
		case ch := <-b.reg:
			b.registers[ch.Chan] = struct{}{}
			ch.Done <- struct{}{}
		//Delete a subscriber.
		case ch, ok := <-b.unreg:
			if ok {
				delete(b.registers, ch.Chan)
				ch.Done <- struct{}{}
			} else {
				return //Terminate this goroutine if this broadcaster is closed.
			}
		}
	}
}
func (b *broadcaster[T]) Register(ch chan<- T) {
	done := make(chan struct{})
	b.reg <- Operate[T]{
		Chan: ch,
		Done: done,
	}
	<-done
}
func (b *broadcaster[T]) Unregister(ch chan<- T) {
	done := make(chan struct{})
	b.unreg <- Operate[T]{
		Chan: ch,
		Done: done,
	}
	<-done
}
func (b *broadcaster[T]) Subbmit(m T) {
	b.input <- m
}
func (b *broadcaster[T]) Close() {
	close(b.input)
	close(b.reg)
	close(b.unreg)
}

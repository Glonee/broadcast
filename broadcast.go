package broadcast

type broadcaster[T any] struct {
	input     chan T
	reg       chan chan<- T
	unreg     chan chan<- T
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
	//return current num of subscribers.
	NumofSubscribers() int
}

//Create a new broadcaster with the given input channel buffer length.
func NewbroadCaster[T any](buflen int) BroadCaster[T] {
	b := &broadcaster[T]{
		input:     make(chan T, buflen),
		reg:       make(chan chan<- T),
		unreg:     make(chan chan<- T),
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
			b.registers[ch] = struct{}{}
		//Delete a subscriber.
		case ch, ok := <-b.unreg:
			if ok {
				delete(b.registers, ch)
			} else {
				return //Terminate this goroutine if this broadcaster is closed.
			}
		}
	}
}
func (b *broadcaster[T]) Register(ch chan<- T) {
	b.reg <- ch
}
func (b *broadcaster[T]) Unregister(ch chan<- T) {
	b.unreg <- ch
}
func (b *broadcaster[T]) Subbmit(m T) {
	b.input <- m
}
func (b *broadcaster[T]) Close() {
	close(b.input)
	close(b.reg)
	close(b.unreg)
}
func (b *broadcaster[T]) NumofSubscribers() int {
	return len(b.registers)
}

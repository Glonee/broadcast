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

//Create a new broadcaster with the given input channel buffer length.
//This broadcaster won't be blocked if unable to send message to given channel.
//Note that the given channel need to have buffer, otherwise it will be ignored.
func NewUnblockedbroadCaster[T any](buflen int) BroadCaster[T] {
	b := &broadcaster[T]{
		input:     make(chan T, buflen),
		reg:       make(chan chan<- T),
		unreg:     make(chan chan<- T),
		registers: make(map[chan<- T]struct{}),
	}
	go b.unblockedrun()
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
func (b *broadcaster[T]) unblockedrun() {
	for {
		select {
		//Send message to all subscribers.
		case m := <-b.input:
			for ch := range b.registers {
				select {
				//If able to send message to this channel
				case ch <- m:
				//Else, ignore this channel
				default:
				}
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

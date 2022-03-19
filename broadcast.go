package broadcast

type broadcaster[T any] struct {
	input     chan T
	reg       chan chan<- T
	unreg     chan chan<- T
	registers map[chan<- T]struct{}
}
type BroadCaster[T any] interface {
	Register(chan<- T)
	Unregister(chan<- T)
	Subbmit(T)
	Close()
}

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
		case m := <-b.input:
			for ch := range b.registers {
				ch <- m
			}
		case ch := <-b.reg:
			b.registers[ch] = struct{}{}
		case ch, ok := <-b.unreg:
			if ok {
				delete(b.registers, ch)
			} else {
				return
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

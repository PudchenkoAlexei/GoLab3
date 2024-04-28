package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	Mq MessageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})
	go func() {
		for !l.stopReq || !l.Mq.empty() {
			op := l.Mq.pull()
			update := op.Do(l.next)
			if update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
}

type MessageQueue struct {
	Ops     []Operation
	mu      sync.Mutex
	blocked chan struct{}
}

func (mq *MessageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.Ops = append(mq.Ops, op)
	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *MessageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for len(mq.Ops) == 0 {
		mq.blocked = make(chan struct{})
		mq.mu.Unlock()
		<-mq.blocked
		mq.mu.Lock()
	}
	op := mq.Ops[0]
	mq.Ops[0] = nil
	mq.Ops = mq.Ops[1:]
	return op
}

func (mq *MessageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.Ops) == 0
}

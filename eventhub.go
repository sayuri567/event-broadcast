package broadcast

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	eventHubs = map[string]*EventHub{}
)

// EventHub EventHub
type EventHub struct {
	sync.Mutex
	eventHandlers []*eventHandler
	quit          chan struct{}
}

// Handler Handler
type Handler interface {
	Handle(msg *Message)
}

// Message Message
type Message struct {
	Type    string
	Message interface{}
}

// NewBroadcaster Add new broadcaster
func NewBroadcaster(name string) (*EventHub, error) {
	if _, ok := eventHubs[name]; ok {
		return nil, fmt.Errorf("broadcaster name %v existed", name)
	}
	eventHubs[name] = &EventHub{
		quit: make(chan struct{}),
	}
	return eventHubs[name], nil
}

// GetBroadcaster Get broadcaster by name, if name not existed, will be return nil
func GetBroadcaster(name string) *EventHub {
	eh, ok := eventHubs[name]
	if !ok {
		return nil
	}
	return eh
}

// AddHandle add event handler
func (slice *EventHub) AddHandle(name string, h Handler) {
	slice.Lock()
	defer slice.Unlock()
	eventHandler := &eventHandler{
		name:   name,
		handle: h,
		quit:   slice.quit,
	}
	eventHandler.start()
	slice.eventHandlers = append(slice.eventHandlers, eventHandler)
}

// Send send message
func (slice *EventHub) Send(msg *Message) {
	slice.Lock()
	defer slice.Unlock()

	for _, eventHandler := range slice.eventHandlers {
		eventHandler.source <- msg
	}
}

// Close close event handler
func (slice *EventHub) Close() {
	close(slice.quit)
	time.Sleep(1 * time.Second)
}

// eventHandler eventHandler
type eventHandler struct {
	name   string
	handle Handler
	source chan *Message
	quit   chan struct{}
}

// Start Start
func (w *eventHandler) start() {
	w.source = make(chan *Message)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				recoverPanic(err)
			}
		}()
		for {
			select {
			case msg := <-w.source:
				if w.handle != nil {
					go func() {
						defer func() {
							if err := recover(); err != nil {
								recoverPanic(err)
							}
						}()
						w.handle.Handle(msg)
					}()
				}
			case <-w.quit:
				logger.Infof("close %v event handler", w.name)
				return
			}
		}
	}()
}

func recoverPanic(r interface{}) {
	if r != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		logger.Errorf(err.Error()+"\nstack: ...\n%s", string(buf))
	}
}

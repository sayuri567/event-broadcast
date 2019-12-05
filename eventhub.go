package broadcast

import (
	"fmt"
	"sync"
	"time"
)

var (
	eventHubs = map[string]*eventHub{}
)

type eventHub struct {
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
func NewBroadcaster(name string) error {
	if _, ok := eventHubs[name]; ok {
		return fmt.Errorf("broadcaster name %v existed.", name)
	}
	eventHubs[name] = &eventHub{
		quit: make(chan struct{}),
	}
	return nil
}

// GetBroadcaster Get broadcaster by name, if name not existed, will be return nil
func GetBroadcaster(name string) *eventHub {
	eh, ok := eventHubs[name]
	if !ok {
		return nil
	}
	return eh
}

func (slice *eventHub) AddHandle(name string, h Handler) {
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

func (slice *eventHub) Iter(msg *Message) {
	slice.Lock()
	defer slice.Unlock()

	for _, eventHandler := range slice.eventHandlers {
		eventHandler.source <- msg
	}
}

func (slice *eventHub) Close() {
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
				logger.Error(err)
			}
		}()
		for {
			select {
			case msg := <-w.source:
				if w.handle != nil {
					go func() {
						defer func() {
							if err := recover(); err != nil {
								logger.Error(err)
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

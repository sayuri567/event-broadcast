package main

import (
	"fmt"

	broadcast "github.com/sayuri567/event-broadcast"
	"github.com/sirupsen/logrus"
)

func main() {
	eh, err := broadcast.NewBroadcaster("test-event")
	if err != nil {
		logrus.WithError(err).Error("failed to new broadcaster")
		return
	}

	h1 := &handler1{}
	h2 := &handler2{}

	eh.AddHandle("handle-1", h1)
	eh.AddHandle("handle-2", h2)

	eh.Send(&broadcast.Message{Type: "test_msg", Message: "event message"})
	eh.Send(&broadcast.Message{Type: "test_msg2", Message: "event message2"})

	eh.Close()
}

type handler1 struct {
}

type handler2 struct {
}

func (h *handler1) Handle(msg *broadcast.Message) {
	if msg.Type != "test_msg" {
		return
	}
	fmt.Println("handler1: ", msg.Message)
}

func (h *handler2) Handle(msg *broadcast.Message) {
	fmt.Println("handler2: ", msg.Message)
}

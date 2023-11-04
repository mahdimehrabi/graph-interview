package message

import (
	"errors"
	"github.com/google/uuid"
	destinationSocket "github.com/mahdimehrabi/graph-interview/broker/external/destination"
	"github.com/mahdimehrabi/graph-interview/broker/external/utils"
	"github.com/mahdimehrabi/graph-interview/broker/internal/entity"
	"github.com/mahdimehrabi/graph-interview/broker/internal/repository/message"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	workerCount          = 50
	saveMessageMethod    = "save_message"
	saveDeadlineDuration = 5 * time.Second
)

var ErrResourceNotAvailable = errors.New("resource is not available")

type destination struct {
	queue   chan *entity.Message
	sockets []*destinationSocket.Socket
}

func NewDestination(sockets []*destinationSocket.Socket) message.Message {
	d := &destination{
		sockets: sockets,
		queue:   make(chan *entity.Message, 10000),
	}
	go d.SaveQueue()
	return d
}

func (b destination) SaveQueue() {
	for i := 0; i < workerCount; i++ {
		go b.savingWorker()
	}
}

func (b destination) savingWorker() {
	for {
		msg := <-b.queue
		id := uuid.New().String()
		socket := b.sockets[utils.RandomNumber(len(b.sockets)-1)]

		deadline := time.NewTicker(saveDeadlineDuration)
		done := make(chan bool)
		go func(ch chan bool) {
			if _, err := socket.SendWaitJSON(msg, saveMessageMethod, id); err != nil {
				log.Printf("failed to save message %s trying again,err:%s", id, err.Error())
				return
			}
			done <- true
		}(done)
		select {
		case <-done:
			log.Printf("message %s saved succesfuly🥳 \n", msg.Message)
		case <-deadline.C: //deadline exceeded
			time.Sleep(1 * time.Microsecond) // socket resend cool down
			b.queue <- msg
		}
	}
}

func (b destination) Save(msg *entity.Message) error {
	if len(b.queue) >= 10000 {
		return ErrResourceNotAvailable
	}
	b.queue <- msg
	return nil
}

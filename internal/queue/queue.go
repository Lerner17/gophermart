package queue

import (
	"errors"

	"github.com/Lerner17/gophermart/internal/models"
)

// In real world we would have Redis or Kafka, etc.
var ordersQueue = make(chan models.OrderMessage, 100000)

var ErrQueueClosed = errors.New("queue is closed")

func PushOrderMessage(msg models.OrderMessage) {
	ordersQueue <- msg
}

func GetNextOrderMessage() (models.OrderMessage, error) {
	msg, ok := <-ordersQueue
	if !ok {
		return msg, ErrQueueClosed
	}
	return msg, nil
}

func DumpAndCloseOrderQueue() []models.OrderMessage {
	close(ordersQueue)
	var result = make([]models.OrderMessage, 0, len(ordersQueue))
	for msg := range ordersQueue {
		result = append(result, msg)
	}
	return result
}

func FullfilQueue(messages []models.OrderMessage) {
	for _, msg := range messages {
		ordersQueue <- msg
	}
}

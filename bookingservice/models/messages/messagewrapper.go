package messages

import "time"

type MessageWrapper[T any] struct {
	ID       string
	Type     string
	CreateAt time.Time
	Payload  T
}

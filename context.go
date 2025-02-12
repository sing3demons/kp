package kp

import "time"

type IContext interface {
	Log(message string)
	Param(name string) string
	Query(name string) string
	ReadInput(data any) error
	Response(responseCode int, responseData any) error

	SendMessage(topic string, message any, opts ...OptionProducerMessage) (RecordMetadata, error)
}

type HandleFunc func(ctx IContext)

type ServiceHandleFunc HandleFunc

type Middleware func(HandleFunc) HandleFunc

type OptionProducerMessage struct {
	key       string
	headers   []map[string]string
	Timestamp time.Time
	Metadata  any
	Offset    int64
	Partition int32
}

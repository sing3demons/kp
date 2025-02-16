package kp

import "context"

type IContext interface {
	Context() context.Context

	Log() ILogger
	Param(name string) string
	Query(name string) string
	ReadInput(data any) error
	Response(code int, data any) error

	SendMessage(topic string, payload any, opts ...OptionProducerMsg) (RecordMetadata, error)
}

type HandleFunc func(ctx IContext) error

type ServiceHandleFunc HandleFunc

type Middleware func(HandleFunc) HandleFunc

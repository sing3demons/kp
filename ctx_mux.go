package kp

import (
	"context"
	"encoding/json"
	"net/http"
)

type HttpContext struct {
	w   http.ResponseWriter
	r   *http.Request
	cfg *KafkaConfig
	log ILogger
}

func newMuxContext(w http.ResponseWriter, r *http.Request, cfg *KafkaConfig, log ILogger) IContext {
	ctx := InitSession(r.Context(), log)
	r = r.WithContext(ctx)

	return &HttpContext{
		w:   w,
		r:   r,
		cfg: cfg,
		log: log,
	}
}

func (c *HttpContext) Context() context.Context {
	return c.r.Context()
}

func (c *HttpContext) SendMessage(topic string, message any, opts ...OptionProducerMsg) (RecordMetadata, error) {
	return producer(c.cfg.producer, topic, message, opts...)
}

func (c *HttpContext) Log() ILogger {
	return c.log
}

func (c *HttpContext) Query(name string) string {
	return c.r.URL.Query().Get(name)
}

func (c *HttpContext) Param(name string) string {
	v := c.r.Context().Value(ContextKey(name))
	var value string
	switch v := v.(type) {
	case string:
		value = v
	}
	c.r = c.r.WithContext(context.WithValue(c.r.Context(), ContextKey(name), nil))
	return value
}

func (c *HttpContext) ReadInput(data any) error {
	return json.NewDecoder(c.r.Body).Decode(data)
}

func (c *HttpContext) Response(responseCode int, responseData any) error {
	c.w.Header().Set("Content-type", "application/json; charset=UTF8")

	c.w.WriteHeader(responseCode)

	err := json.NewEncoder(c.w).Encode(responseData)
	return err
}

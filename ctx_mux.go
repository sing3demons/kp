package kp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpContext struct {
	w   http.ResponseWriter
	r   *http.Request
	cfg *KafkaConfig
}

func newMuxContext(w http.ResponseWriter, r *http.Request, cfg *KafkaConfig) IContext {
	return &HttpContext{
		w:   w,
		r:   r,
		cfg: cfg,
	}
}

func (c *HttpContext) SendMessage(topic string, message interface{}, opts ...OptionProducerMessage) (RecordMetadata, error) {
	return RecordMetadata{}, nil
}

func (c *HttpContext) Log(message string) {
	fmt.Println("Context:", message)
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

func (c *HttpContext) ReadInput(data interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(data)
}

func (c *HttpContext) Response(responseCode int, responseData interface{}) error {
	c.w.Header().Set("Content-type", "application/json; charset=UTF8")

	c.w.WriteHeader(responseCode)

	return json.NewEncoder(c.w).Encode(responseData)
}

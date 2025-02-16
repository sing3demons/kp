package kp

import (
	"context"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	ctx *gin.Context
	cfg *KafkaConfig
	log ILogger
}

func newGinContext(c *gin.Context, cfg *KafkaConfig, log ILogger) IContext {
	ctx := InitSession(c.Request.Context(), log)
	c.Request = c.Request.WithContext(ctx)
	return &GinContext{ctx: c, cfg: cfg, log: log}
}

func (c *GinContext) Context() context.Context {
	return c.ctx.Request.Context()
}

func (c *GinContext) SendMessage(topic string, message any, opts ...OptionProducerMsg) (RecordMetadata, error) {
	return producer(c.cfg.producer, topic, message, opts...)
}

func (c *GinContext) Log() ILogger {
	return c.log
}

func (c *GinContext) Query(name string) string {
	return c.ctx.Query(name)
}

func (c *GinContext) Param(name string) string {
	return c.ctx.Param(name)
}

func (c *GinContext) ReadInput(data any) error {
	return c.ctx.BindJSON(data)
}

func (c *GinContext) Response(responseCode int, responseData any) error {
	c.ctx.JSON(responseCode, responseData)
	return nil
}

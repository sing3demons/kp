package kp

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// context consumer
type ConsumerContext struct {
	message string
	msg     *kafka.Message
	cfg     *KafkaConfig
}

func NewConsumerContext(cfg *KafkaConfig, msg *kafka.Message) IContext {
	return &ConsumerContext{
		message: string(msg.Value),
		msg:     msg,
		cfg:     cfg,
	}
}

func (c *ConsumerContext) SendMessage(topic string, message any, opts ...OptionProducerMessage) (RecordMetadata, error) {
	panic("not implemented")
}

func (c *ConsumerContext) Log(message string) {
	log.Println("Context:", message)
}

func (c *ConsumerContext) Query(name string) string {
	return ""
}

func (c *ConsumerContext) Param(name string) string {
	return ""
}

func (ctx *ConsumerContext) ReadInput(data any) error {
	fmt.Println("ReadInput:", ctx.message)
	const errMsgFormat = "%s, payload: %s"
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if val.Elem().Kind() == reflect.String {
			val.Elem().SetString(ctx.message)
			return nil
		}

		if err := json.Unmarshal([]byte(ctx.message), data); err != nil {
			return fmt.Errorf(errMsgFormat, err.Error(), ctx.message)
		}
		return nil
	case reflect.String:
		return fmt.Errorf("cannot assign to non-pointer string")
	default:
		err := json.Unmarshal([]byte(ctx.message), &data)
		if err != nil {
			return fmt.Errorf(errMsgFormat, err.Error(), ctx.message)
		}
		return nil
	}
}

func (c *ConsumerContext) Response(responseCode int, responseData any) error {
	fmt.Println("Response:", responseData)
	return nil
}

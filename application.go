package kp

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IApplication interface {
	Get(path string, handler HandleFunc, middlewares ...Middleware)
	Post(path string, handler HandleFunc, middlewares ...Middleware)
	Put(path string, handler HandleFunc, middlewares ...Middleware)
	Patch(path string, handler HandleFunc, middlewares ...Middleware)
	Delete(path string, handler HandleFunc, middlewares ...Middleware)
	Use(middlewares ...Middleware)
	Start()

	Consume(topic string, h ServiceHandleFunc) error
}

type AppConfig struct {
	Port   string
	Router Router
}

type KafkaConfig struct {
	Brokers     string
	GroupID     string
	exitChannel chan bool

	AutoCreateTopic bool
}

type Config struct {
	AppConfig   AppConfig
	KafkaConfig KafkaConfig
	exitChannel chan bool
}

// enum Router {gin, mux}
type Router int

const (
	None Router = iota
	Gin
	Mux
	Fiber
	Kafka
)

func NewApplication(cfg Config) IApplication {
	switch cfg.AppConfig.Router {
	case Gin:
		return newGinServer(cfg)
	case Mux:
		return newServer(cfg)
	case Kafka:
		return newKafkaConsumer(cfg)
	default:
		return newServer(cfg)
	}
}

type consumerApplication struct {
	cfg *Config
}

func newKafkaConsumer(cfg Config) IApplication {
	return &consumerApplication{
		cfg: &cfg,
	}
}

func (ms *consumerApplication) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	panic("not implemented")
}

func (ms *consumerApplication) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	panic("not implemented")
}

func (ms *consumerApplication) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	panic("not implemented")
}

func (ms *consumerApplication) Patch(path string, handler HandleFunc, middlewares ...Middleware) {
	panic("not implemented")
}

func (ms *consumerApplication) Delete(path string, handler HandleFunc, middlewares ...Middleware) {
	panic("not implemented")
}

func (ms *consumerApplication) Use(middlewares ...Middleware) {
	panic("not implemented")
}

func newConsumer(cfg Config) *consumerApplication {
	return &consumerApplication{
		cfg: &cfg,
	}
}

func (ms *consumerApplication) newKafkaConsumer(servers string, groupID string) (*kafka.Consumer, error) {
	// Configurations
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	config := &kafka.ConfigMap{

		// Alias for metadata.broker.list: Initial list of brokers as a CSV list of broker host or host:port.
		// The application may also use rd_kafka_brokers_add() to add brokers during runtime.
		"bootstrap.servers": servers,

		// Client group id string. All clients sharing the same group.id belong to the same group.
		"group.id": groupID,

		// Action to take when there is no initial offset in offset store or the desired offset is out of range:
		// 'smallest','earliest' - automatically reset the offset to the smallest offset,
		// 'largest','latest' - automatically reset the offset to the largest offset,
		// 'error' - trigger an error which is retrieved by consuming messages and checking 'message->err'.
		"auto.offset.reset": "earliest",

		// Protocol used to communicate with brokers.
		// plaintext, ssl, sasl_plaintext, sasl_ssl
		"security.protocol": "plaintext",

		// Automatically and periodically commit offsets in the background.
		// Note: setting this to false does not prevent the consumer from fetching previously committed start offsets.
		// To circumvent this behaviour set specific start offsets per partition in the call to assign().
		"enable.auto.commit": true,

		// The frequency in milliseconds that the consumer offsets are committed (written) to offset storage. (0 = disable).
		// default = 5000ms (5s)
		// 5s is too large, it might cause double process message easily, so we reduce this to 200ms (if we turn on enable.auto.commit)
		"auto.commit.interval.ms": 500,

		// Automatically store offset of last message provided to application.
		// The offset store is an in-memory store of the next offset to (auto-)commit for each partition
		// and cs.Commit() <- offset-less commit
		"enable.auto.offset.store": true,

		// Enable TCP keep-alives (SO_KEEPALIVE) on broker sockets
		"socket.keepalive.enable": true,
	}

	kc, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return kc, err
}

func (ms *consumerApplication) consumeSingle(topic string, h ServiceHandleFunc) {

	var readTimeout time.Duration = -1
	c, err := ms.newKafkaConsumer(ms.cfg.KafkaConfig.Brokers, ms.cfg.KafkaConfig.GroupID)
	if err != nil {
		return
	}

	defer c.Close()

	if ms.cfg.KafkaConfig.AutoCreateTopic {
		kafkaAdmin, err := kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers": ms.cfg.KafkaConfig.Brokers,
		})

		if err == nil {
			// Create topic
			topicSpec := kafka.TopicSpecification{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			}

			_, err = kafkaAdmin.CreateTopics(context.Background(), []kafka.TopicSpecification{topicSpec})
		}

	}

	c.Subscribe(topic, nil)

	for {
		if readTimeout <= 0 {
			// readtimeout -1 indicates no timeout
			readTimeout = -1
		}

		msg, err := c.ReadMessage(readTimeout)
		if err != nil {
			kafkaErr, ok := err.(kafka.Error)
			if ok {
				if kafkaErr.Code() == kafka.ErrTimedOut {
					if readTimeout == -1 {
						// No timeout just continue to read message again
						continue
					}
				}
			}
			return
		}

		// Execute Handler
		h(NewConsumerContext(&ms.cfg.KafkaConfig, msg))
	}
}

func (ms *consumerApplication) Consume(topic string, h ServiceHandleFunc) error {
	go ms.consumeSingle(topic, h)
	return nil
}

func (ms *consumerApplication) Start() {
	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	ms.cfg.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			exit = true
		case <-ms.cfg.exitChannel:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			exit = true
		}
	}

	return
}

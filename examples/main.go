package main

import (
	"fmt"

	"github.com/sing3demons/kp"
)

func main() {
	servers := "localhost:29092"
	topic := "when-citizen-has-registered"
	groupID := "validation-consumer"

	ms := kp.NewApplication(kp.Config{
		AppConfig: kp.AppConfig{
			Router: kp.Kafka,
		},
		KafkaConfig: kp.KafkaConfig{
			Brokers:         servers,
			GroupID:         groupID,
			AutoCreateTopic: true,
		},
	})

	ms.Consume(topic, func(ctx kp.IContext) {
		fmt.Println("Message: ")

		var message string
		ctx.ReadInput(&message)
		fmt.Println("Message: ", message)
	})

	ms.Consume("order", func(ctx kp.IContext) {
		fmt.Println("Order: ")

		var order string
		ctx.ReadInput(&order)
		fmt.Println("Order: ", order)
	})

	// defer ms.Cleanup()
	ms.Start()
}

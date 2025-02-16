package main

import (
	"fmt"
	"net/http"

	"github.com/sing3demons/kp"
)

func main() {
	servers := "localhost:29092"
	topic := "test"
	groupID := "validation-consumer"

	logger := kp.NewZapLogger(kp.NewAppLogger())
	defer logger.Sync()

	ms := kp.NewApplication(&kp.Config{
		AppConfig: kp.AppConfig{
			Port: "3000",
		},
		KafkaConfig: kp.KafkaConfig{
			Brokers: []string{servers},
			GroupID: groupID,
		},
	}, logger)

	ms.Get("/", func(ctx kp.IContext) error {
		log := ctx.Log().L(ctx.Context())
		name := ctx.Param("name")

		log.Info(fmt.Sprintf("Request with name: %s", name))

		if name == "test" {
			ctx.SendMessage("test", "Hello, Test!")
			return ctx.Response(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
		}

		ctx.SendMessage("order", "Hello, World!")
		return ctx.Response(http.StatusOK, "Hello, World!")
	})

	ms.Consume(topic, func(ctx kp.IContext) error {
		fmt.Println("Message: ")

		var message string
		ctx.ReadInput(&message)
		fmt.Println("Message: ", message)

		return nil
	})

	ms.Consume("order", func(ctx kp.IContext) error {
		log := ctx.Log().L(ctx.Context())
		var order string
		ctx.ReadInput(&order)
		log.Info(fmt.Sprintf("testmessage: %v", order))

		log.Info("========== test===========")
		return nil
	})

	// defer ms.Cleanup()
	ms.Start()
}

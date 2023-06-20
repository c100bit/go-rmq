package main

import (
	"context"
	"eventdrivenrabbit/internal"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("guest", "guest", "localhost:5672", "/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}
	var blocking chan struct{}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Printf("new msg %v", msg)
				time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Printf("ack")
					return err
				}
				log.Printf("Ack msg %s", message.MessageId)
				return nil
			})
		}
	}()
	log.Println("")
	<-blocking
}

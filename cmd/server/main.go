package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
		return
	}

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Pausing the game...")
			sendMessage(ch, "pause")
		case "resume":
			fmt.Println("Resuming the game...")
			sendMessage(ch, "resume")
		case "quit":
			fmt.Println("Exiting the game...")
			return
		default:
			fmt.Println("Command not found...")
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()
	fmt.Println("Received shutdown signal, shutting down...")
}

func sendMessage(ch *amqp.Channel, message string) {
	var isPaused bool
	if message == "pause" {
		isPaused = true
	} else if message == "resume" {
		isPaused = false
	}

	err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: isPaused,
	})
	if err != nil {
		fmt.Println("Failed to publish a message:", err)
	}
}

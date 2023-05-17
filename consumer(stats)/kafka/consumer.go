package kafka

import (
	"log"
	"fmt"
	"context"

	"github.com/Shopify/sarama"
)

var (
	kafkaBrokers = []string{"localhost:9093"}
	kafkaTopics = []string{"wallet_operations"}
	consumerGroupID = "sarama_consumer"
)

func StartConsumer() {
	config := sarama.NewConfig()

	fmt.Println("HELLO KAFKA")

	client, err := sarama.NewClient(kafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	defer func () { _ = client.Close() }()

	group, err := sarama.NewConsumerGroupFromClient(consumerGroupID, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()
	log.Println("Kafka __consumer__ running!")

	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	ctx := context.Background()
	for {
		handler := ConsumerGroupHandler{}

		err := group.Consume(ctx, kafkaTopics, handler)
		if err != nil {
			panic(err)
		}
	}
}
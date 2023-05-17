package kafka

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

var (
	KafkaBrokers = []string{"localhost:9093"}
	KafkaTopics = "wallet_operations"
	enqueued int
)

func PublishEvent(event interface{}) {

	println("Publish event function")
	producer, err := setupProducer()
	if err != nil {
		panic(err)
	} else {
		log.Printf("Kafka __producer__ running!")
	}

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	produceMessage(producer, signals, event)
}

func setupProducer() (sarama.AsyncProducer, error) {
	
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(KafkaBrokers, config)
}

func produceMessage(producer sarama.AsyncProducer, signals chan os.Signal, event interface{}) {
	message := &sarama.ProducerMessage{
        Topic: KafkaTopics,
		Key:	sarama.StringEncoder("def_key"), 
		Value: 	sarama.StringEncoder("def_value"),
    }
    switch event := event.(type) {
    	case WalletCreatedEvent:
       		message.Value = sarama.StringEncoder(fmt.Sprintf(`{"type":"%s", "id":"%s"}`, event.Type, event.Id))
			message.Key = sarama.StringEncoder(event.Type)
		case WalletDeletedEvent:
			message.Value = sarama.StringEncoder(fmt.Sprintf(`{"type":"%s", "id":"%s", "status":"%s"}`, event.Type, event.Id, event.Status))
			message.Key = sarama.StringEncoder(event.Type)
		case WalletDepositedEvent:
			message.Value = sarama.StringEncoder(fmt.Sprintf(`{"type":"%s", "amount":"%s"}`, event.Type, event.Amount))
			message.Key = sarama.StringEncoder(event.Type)
		case WalletWithdrawnEvent:
			message.Value = sarama.StringEncoder(fmt.Sprintf(`{"type":"%s", "amount":"%s"}`, event.Type, event.Amount))
			message.Key = sarama.StringEncoder(event.Type)
		case WalletTransferedEvent:
			message.Value = sarama.StringEncoder(fmt.Sprintf(`{"type":"%s", "amount":"%s"}`, event.Type, event.Amount))
			message.Key = sarama.StringEncoder(event.Type)
    default:
        log.Fatalln("Unsupported event type")
    }

	select {
		case producer.Input() <- message:
			enqueued++
			log.Println("New event created and received")
		case <- signals:
			producer.AsyncClose()
			return
	}
}
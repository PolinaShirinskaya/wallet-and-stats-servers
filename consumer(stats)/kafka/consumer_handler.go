package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"stats_wallets.com/stats"
)

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d message: %v\n",
			message.Topic, message.Partition, message.Offset, string(message.Value))

		sess.MarkMessage(message, "")

		eventType := stats.DetermineEventType(message.Key)

		switch eventType {
		case "Wallet_Created":
			var createdEvent stats.WalletCreatedEvent
			err := json.Unmarshal(message.Value, &createdEvent)
			if err != nil {
				log.Printf("Error deserializing Wallet_Created event: %s\n", err.Error())
				continue
			}
			stats.CreatingUpdate(createdEvent)
			log.Printf("Handling Wallet_Created event: %+v\n", createdEvent)

		case "Wallet_Deleted":
			var deletedEvent stats.WalletDeletedEvent
			err := json.Unmarshal(message.Value, &deletedEvent)
			if err != nil {
				log.Printf("Error deserializing Wallet_Deleted event: %s\n", err.Error())
				continue
			}
			stats.DeletingUpdate(deletedEvent)
			log.Printf("Handling Wallet_Deleted event: %+v\n", deletedEvent)

		case "Wallet_Deposited":
			var depositedEvent stats.WalletDepositedEvent
			err := json.Unmarshal(message.Value, &depositedEvent)
			if err != nil {
				log.Printf("Error deserializing Wallet_Deposited event: %s\n", err.Error())
				continue
			}
			stats.DepositingUpdate(depositedEvent)
			log.Printf("Handling Wallet_Deposited event: %+v\n", depositedEvent)

		case "Wallet_Withdrawn":
			var wihdrawnEvent stats.WalletWithdrawnEvent
			err := json.Unmarshal(message.Value, &wihdrawnEvent)
			if err != nil {
				log.Printf("Error deserializing Wallet_Withdrawn event: %s\n", err.Error())
				continue
			}
			stats.WithdrawingUpdate(wihdrawnEvent)
			log.Printf("Handling Wallet_Wihtdarwn event: %+v\n", wihdrawnEvent)

		case "Wallet_Transfered":
			var transferesEvent stats.WalletTransferedEvent
			err := json.Unmarshal(message.Value, &transferesEvent)
			if err != nil {
				log.Printf("Error deserializing Wallet_Deposited event: %s\n", err.Error())
				continue
			}
			stats.TransfertingUpdate(transferesEvent)
			log.Printf("Handling Wallet_Deposited event: %+v\n", transferesEvent)

		default:
			log.Printf("Unknown event")
		}
	}
	return nil
}

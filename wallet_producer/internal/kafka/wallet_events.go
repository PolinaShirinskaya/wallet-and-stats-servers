package kafka

import "strconv"

//"encoding/json"

type WalletCreatedEvent struct {
	Type 	string `json:"type"`
	Id		string `json:"id"`
	Status	string `json:"status"`
}

type WalletDeletedEvent struct {
	Type 	string `json:"type"`
	Id		string `json:"id"`
	Status	string `json:"status"`
}

type WalletDepositedEvent struct {
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}

type WalletWithdrawnEvent struct {
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}

type WalletTransferedEvent struct {
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}


func WalletСreateEvent(id string, status string) {
	event := WalletCreatedEvent{
		Type: "Wallet_Created",
		Id: id,
		Status: status,
	}
	println("Wallet create event function")
	PublishEvent(event)
}

func WalletDeleteEvent(id string, status string) {
	event := WalletDeletedEvent{
		Type: "Wallet_Deleted",
		Id: id,
		Status: status,
	}
	PublishEvent(event)
}

func WalletDepositeEvent(amount float64) {
	convAmount := strconv.FormatFloat(amount, 'f', -1, 64)
	event := WalletDepositedEvent{
		Type: "Wallet_Deposited",
		Amount: convAmount,
	}
	PublishEvent(event)
}

func WalletWithdrawEvent(amount float64) {
	convAmount := strconv.FormatFloat(amount, 'f', -1, 64)
	event := WalletWithdrawnEvent{
		Type: "Wallet_Withdrawn",
		Amount: convAmount,
	}
	PublishEvent(event)
}

func WalletTransferEvent(amount float64) {
	convAmount := strconv.FormatFloat(amount, 'f', -1, 64)
	event := WalletTransferedEvent{
		Type: "Wallet_Transfered",
		Amount: convAmount,
	}
	PublishEvent(event)
}

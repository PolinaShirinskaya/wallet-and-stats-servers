package stats

import (
	"sync"
	"fmt"

	"github.com/shopspring/decimal"

)

var Lock sync.Mutex

type WalletStats struct {
	Wallet struct {
		Total      int     `json:"total"`
		Active     int     `json:"active"`
		Inactive   int     `json:"inactive"`
		Deposited  float64 `json:"deposited"`
		Withdrawn  float64 `json:"withdrawn"`
		Transferred float64 `json:"transferred"`
	} `json:"wallet"`
}


var ResultStats WalletStats

func CreatingUpdate(event WalletCreatedEvent) {
	Lock.Lock()
	defer Lock.Unlock()


	ResultStats.Wallet.Total++
	ResultStats.Wallet.Active++
}

func DeletingUpdate(event WalletDeletedEvent) {
	Lock.Lock()
	defer Lock.Unlock()

	ResultStats.Wallet.Inactive++
	ResultStats.Wallet.Active--
}

func DepositingUpdate(event WalletDepositedEvent) {
	Lock.Lock()
	defer Lock.Unlock()

	
	amount, err := decimal.NewFromString(event.Amount)
	if err != nil {
		fmt.Println("Decimal error: convert string to float64", err)
		return
	}

	deposit := decimal.NewFromFloat(ResultStats.Wallet.Deposited)
	sum := deposit.Add(amount)

	ResultStats.Wallet.Deposited, _ = sum.Float64()

}

func WithdrawingUpdate(event WalletWithdrawnEvent) {
	Lock.Lock()
	defer Lock.Unlock()

	amount, err := decimal.NewFromString(event.Amount)
	if err != nil {
		fmt.Println("Decimal error: convert string to float64", err)
		return
	}

	deposit := decimal.NewFromFloat(ResultStats.Wallet.Withdrawn)
	sum := deposit.Sub(amount)

	ResultStats.Wallet.Withdrawn, _ = sum.Float64()
}

func TransfertingUpdate(event WalletTransferedEvent) {
	Lock.Lock()
	defer Lock.Unlock()

	amount, err := decimal.NewFromString(event.Amount)
	if err != nil {
		fmt.Println("DEcimal error: convert string to float64", err)
		return
	}

	deposit := decimal.NewFromFloat(ResultStats.Wallet.Transferred)
	sum := deposit.Add(amount)

	ResultStats.Wallet.Transferred, _ = sum.Float64()
}
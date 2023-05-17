package walletstore

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	Id string 		`json:"id"`
	Name string 	`json:"name"`
	Status string 	`json:"status"`
	Balance float64 `json:"balance"`
}

type WalletStore struct {
	sync.Mutex

	wallets []Wallet
}

func New() *WalletStore {
	ws := &WalletStore{}
	ws.wallets = make([]Wallet, 0)
	return ws
}

//для создания собственного типа ошибок (использую в TransferWallet function)
type MyError struct {
	Message string
}

func (e MyError) Error() string {
	return fmt.Sprintf(e.Message)
}

func addFloat64(balance, amount float64) float64 {
	decBalance := decimal.NewFromFloat(balance)
	decAmount := decimal.NewFromFloat(amount)
	result := decBalance.Add(decAmount)
	res, _ := result.Float64()
	return res
}

func subtractFloat64(balance, amount float64) float64 {
	decBalance := decimal.NewFromFloat(balance)
	decAmount := decimal.NewFromFloat(amount)
	result := decBalance.Sub(decAmount)
	res, _ := result.Float64()
	return res
}

//создание коешлька
func (ws *WalletStore) CreateWallet(name string) (Wallet, error) {
	ws.Lock()
	defer ws.Unlock()
	
	if len(name) > 1 {
		var wallet Wallet
		wallet.Id = strconv.Itoa(rand.Intn(100000000))
		wallet.Name = name
		wallet.Status = "active"
		ws.wallets = append(ws.wallets, wallet)
		return wallet, nil
	} else {
		return Wallet{}, MyError{Message: "'Name' must consist of more than one symbol"}
	}	
}

//получение списка всех кошельков(в т.ч. и деактивированных)
func (ws *WalletStore) GetAllWallets() []Wallet {
	ws.Lock()
	defer ws.Unlock()

	return ws.wallets
}

//получение кошелька по его ID
func (ws *WalletStore) GetWallet(id string) (Wallet, error) {
	ws.Lock()
	defer ws.Unlock()

	for _, wallet := range ws.wallets {
		if wallet.Id == id {
			return wallet, nil
		}
	}
	return Wallet{}, MyError{"Wallet not found"}
}

//деактивация кошелька(status ---> inactive)
func (ws *WalletStore) DeleteWallet(id string) (Wallet, error) {
	ws.Lock()
	defer ws.Unlock()

	for index, wallet := range ws.wallets {
		if wallet.Id == id {
			ws.wallets = append(ws.wallets[:index], ws.wallets[index+1:]...)
			wallet.Status = "inactive"
			ws.wallets = append(ws.wallets, wallet)
			return wallet, nil
		}
	}
	return Wallet{}, MyError{"Wallet not found"}	
}

//обновление имени кошелька
func (ws *WalletStore) UpdateWallet(name, id string) error {
	ws.Lock()
	defer ws.Unlock()

	if len(name) <= 1 {
		return MyError{"'Name' must consist of more than one symbol"}
	} else {

		for index, wallet := range ws.wallets {
			if wallet.Id == id && wallet.Status == "active"{
				ws.wallets = append(ws.wallets[:index], ws.wallets[index+1:]...)
				wallet.Name = name
				ws.wallets = append(ws.wallets, wallet)
				return nil
			}
		}
	}
	return MyError{"Wallet not found"}
}

//пополнение средств
func (ws *WalletStore) DepositWallet(amount float64, id string) error {
	ws.Lock()
	defer ws.Unlock()

	for index, wallet := range ws.wallets {
		if wallet.Id == id && wallet.Status == "active" {
			if amount >= 0 {
				ws.wallets = append(ws.wallets[:index], ws.wallets[index+1:]...)
				wallet.Balance = addFloat64(wallet.Balance, amount)
				ws.wallets = append(ws.wallets, wallet)
				return nil
			} else {
				return MyError{Message: "Amount cannot be negative"}
			}
		}
	}
	return MyError{Message: "Wallet not found"}
}

//снятие средств
func (ws *WalletStore) WithdrawWallet(amount float64, id string) error {
	ws.Lock()
	defer ws.Unlock()

	for index, wallet := range ws.wallets {
		if wallet.Id == id && wallet.Status == "active"{
			if amount > 0 {
				ws.wallets = append(ws.wallets[:index], ws.wallets[index+1:]...)
				wallet.Balance = subtractFloat64(wallet.Balance, amount)
				if wallet.Balance > 0 {
					ws.wallets = append(ws.wallets, wallet)
					return nil
				} else {
					return MyError{Message: "Insufficient funds in the wallet"}
				}
			} else {
				return MyError{Message: "Amount cannot be negative"}
			}
		}
	}
	return MyError{Message: "Wallet not found"}
}

//перевод средств между счетами
func (ws *WalletStore) TransferWallet(amount float64, senderID string, recipientID string) error {
	ws.Lock()
	defer ws.Unlock()

	var senderIndex, recipientIndex int
	var senderWallet, recipientWallet Wallet

	for index, wallet := range ws.wallets {
		if wallet.Id == senderID && wallet.Status == "active" {
			senderIndex = index
			senderWallet = wallet
		} else if wallet.Id == recipientID && wallet.Status == "active" {
			recipientIndex = index
			recipientWallet = wallet
		}
	}

	if senderWallet.Id == "" {
		return MyError{Message: "Sender wallet not found"}
	}
	if recipientWallet.Id == "" {
		return MyError{Message: "Recipient wallet not found"}
	}

	if amount <= 0 {
		return MyError{Message: "Amount cannot be negative"}
	}
	if senderWallet.Balance < amount {
		return MyError{Message: "Insufficient funds in the sender's wallet"}
	}

	senderWallet.Balance = subtractFloat64(senderWallet.Balance, amount)

	recipientWallet.Balance = addFloat64(recipientWallet.Balance, amount)

	ws.wallets[senderIndex] = senderWallet
	ws.wallets[recipientIndex] = recipientWallet

	return nil
}



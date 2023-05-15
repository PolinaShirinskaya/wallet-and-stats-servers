package stats

func DetermineEventType(key []byte) string {
	keyString := string(key)

	switch keyString {
	case "Wallet_Created":
		return "Wallet_Created"
	case "Wallet_Deleted":
		return "Wallet_Deleted"
	case "Wallet_Deposited":
		return "Wallet_Deposited"
	case "Wallet_Withdrawn":
		return "Wallet_Deposited"
	case "Wallet_Transfered":
		return "Wallet_Deposited"
	default:
		return "Unknown"
	}
}

type WalletCreatedEvent struct {
	Type 	string `json:"type"`
	Id		string `json:"id"`
}

type WalletDeletedEvent struct {
	Key    string `json:"-"`
	Type 	string `json:"type"`
	Id		string `json:"id"`
	Status	string `json:"status"`
}

type WalletDepositedEvent struct {
	Key    string `json:"-"`
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}

type WalletWithdrawnEvent struct {
	Key    string `json:"-"`
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}

type WalletTransferedEvent struct {
	Key    string `json:"-"`
	Type	string `json:"type"`
	Amount	string `json:"amount"`
}


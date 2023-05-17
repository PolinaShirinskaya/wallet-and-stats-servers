package walletstore

import (
	"testing"
)

func TestCreate(t *testing.T) {
	ws := New()
	wallet, err := ws.CreateWallet("ValidName")

	//проверка полей созданного кошелька
	if err == nil && wallet.Id == "" {
		t.Error("expected non-empty wallet ID")
	}
	if err == nil && wallet.Name != "ValidName" {
		t.Errorf("expected name 'Test Wallet', got '%s'", wallet.Name)
	}
	if err == nil && wallet.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", wallet.Status)
	}
	

	_, errorName := ws.CreateWallet("V")
	if errorName == nil {
		t.Error("create wallet with invalid name")
	}	
}

func TestGetMethods(t *testing.T) {

	ws := New()
	wallet, _ := ws.CreateWallet("Golang")

	getWallet, err := ws.GetWallet(wallet.Id)
	if err != nil {
		t.Fatal(err)
	}

	// Сравнение ID	созданного кошелька и полученного
	if getWallet.Id != wallet.Id {
		t.Errorf("got wallet.ID = %s, want id = %s", getWallet.Id, wallet.Id)
	}
	// Сравнение NAME созданного и полученного
	if getWallet.Name != "Golang" {
		t.Errorf("got wallet.Name = %s, want name = %s", getWallet.Name, "Golang")
	}
	// Попытка получить кошелек по несуществующему ID
	_, err = ws.GetWallet("123")
	if err == nil {
		t.Fatal("got nill, want error")
	}

	// Запрос списка всех кошельков
	allWallets := ws.GetAllWallets()
	if len(allWallets) != 1 {
		t.Errorf("got len(allWallets) = %d, want 1", len(allWallets))
	}

	// Создадим еще один кошелек. Ожидаем в списке получить уже два
	ws.CreateWallet("Google")
	allWallets2 := ws.GetAllWallets()
	if len(allWallets2) != 2{
		t.Errorf("got len(allWallets2)=%d, want 2", len(allWallets2))
	}
}

func TestDelete(t *testing.T) {
	ws := New()
	delWallet,_ := ws.CreateWallet("First")

	// Проверяем изменился ли status кошелька после деактивации
	ws.DeleteWallet(delWallet.Id)
	getWallet, _ := ws.GetWallet(delWallet.Id)
	if getWallet.Status != "inactive" {
		t.Error("after delete wallet - status should be 'inactive'")
	}

	// Попытка удалить кошелек с несуществующим ID
	if _, err := ws.DeleteWallet("123"); err == nil {
		t.Fatal("try delete wallet with non - existent ID, got no error; want error")
	}
}

func TestUpdate(t *testing.T) {
	ws := New()
	uptWallet,_ := ws.CreateWallet("GoLang")

	// Попытка обновить name кошелька
	ws.UpdateWallet("Yandex", uptWallet.Id)
	getWallet, _ := ws.GetWallet(uptWallet.Id)
	if getWallet.Name != "Yandex" {
		t.Errorf("name of wallet does not update")
	}

	// Попытка обновить name  несуществующего кошелька
	err := ws.UpdateWallet("New", "123")
	if err == nil {
		t.Fatal("got nill, want error: try to update non-existing wallet")
	}
}

func TestDeposit(t *testing.T) {
	ws := New()
	depWallet,_ := ws.CreateWallet("Deposit")

	// Пополнение отрицательной суммой
	error := ws.DepositWallet(-500.21, depWallet.Id)
	if error == nil {
		t.Error("attempt to deposit negative amount")
	}

	// Попытка пополнить несуществующий кошелек
	err := ws.DepositWallet(500.42, "123")
	if err == nil {
		t.Fatal("got nill, want error: try to deposit non-existing wallet")
	}

	// Попытка пополнить деактивированный кошелек
	ws.DeleteWallet(depWallet.Id)
	ws.DepositWallet(500.42, depWallet.Id)
	if depWallet.Status == "inactive" {
		t.Errorf("try to deposit inactive wallet")
	}
}

func TestWithdraw(t *testing.T) {
	ws := New()
	withdrawWallet,_ := ws.CreateWallet("Withdraw")

	// Пополнение отрицательной суммой
	error := ws.WithdrawWallet(-500.21, withdrawWallet.Id)
	if error == nil {
		t.Error("attempt to deposit negative amount")
	}

	// Снятие суммы большей баланса на кошельке
	ws.WithdrawWallet(100, withdrawWallet.Id)
	er := ws.WithdrawWallet(-500.21, withdrawWallet.Id)
	if er == nil {
		t.Error("insufficient funds to withdraw")
	}

	// Попытка пополнить несуществующий кошелек
	err := ws.WithdrawWallet(500.42, "123")
	if err == nil {
		t.Fatal("got nill, want error: try to deposit non-existing wallet")
	}

	// Попытка пополнить деактивированный кошелек
	ws.DeleteWallet(withdrawWallet.Id)
	ws.WithdrawWallet(500.42, withdrawWallet.Id)
	if withdrawWallet.Status == "inactive" {
		t.Errorf("try to deposit inactive wallet")
	}
}

func TestTransfer(t *testing.T) {
	ws := New()
	sender,_ := ws.CreateWallet("Sender")
	recipient,_ := ws.CreateWallet("Recipient")


	// Попытка провести операцию с деактивированными кошелеками
	ws.DepositWallet(1000, sender.Id)
	ws.DeleteWallet(recipient.Id)
	ws.DeleteWallet(sender.Id)
	err := ws.TransferWallet(500.42, sender.Id, recipient.Id)
	if err == nil {
		t.Errorf("operation with inactive wallets")
	}

	// Операция с несуществующими кошельками
	error := ws.TransferWallet(500.42, "123", "321")
	if error == nil {
		t.Errorf("transfer with non - exixtent wallets")
	}

	// Снятие суммы большей имеющейся на кошельке у отправителя
	er := ws.TransferWallet(100000, sender.Id, recipient.Id)
	if er == nil {
		t.Errorf("the sender's balance cannot become negative")
	}

	// Операция с отрицательной суммой перевода
	erro := ws.TransferWallet(-500.42, sender.Id, recipient.Id)
	if erro == nil {
		t.Errorf("try to transfer a negative amount")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
	"os"
	"os/signal"

	"example.com/internal/walletstore"
	"example.com/internal/kafka"
)

func main() {
	mux := http.NewServeMux()
	server := NewWalletServer()
	mux.HandleFunc("/wallets/", server.walletHandler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		log.Println("Received interrupt signal. Shutting down...")
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":5100", mux))
	
}

type walletServer struct {
	store *walletstore.WalletStore
}

// конкструктор для server type, server оборачивает WalletStore,
// который безопасен для concurrent access
func NewWalletServer() *walletServer {
	store := walletstore.New()
	return &walletServer{store: store}
}

// обработка RequestBody
func handlingRequest(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")         //возвращает первое значение, связанное с этим ключом
	mediatype, _, err := mime.ParseMediaType(contentType) //ананлизирует и возвращает преобразованный в нижний регистр и очищенный от пробелов
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

// форматирование Response в JSON формат
func renderJSON(w http.ResponseWriter, v interface{}) {
	json, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ws *walletServer) walletHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/wallets/" {
		//запросы строго "/wallets/" без ID
		if req.Method == http.MethodPost {
			ws.createWalletHandler(w, req)
		} else if req.Method == http.MethodGet {
			ws.getAllWalletsHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /wallets/, got %v", req.Method), http.StatusMethodNotAllowed)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		//запросы имеющие {id}
		path := strings.Trim(req.URL.Path, "/")
		pathParams := strings.Split(path, "/")
		if len(pathParams) < 2 {
			http.Error(w, "expect /wallets/{id} in wallet handler", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
		}
		id := pathParams[1]

		//общие операции с кошельком по его ID
		if len(pathParams) == 2 {
			if req.Method == http.MethodGet {
				ws.getWalletHandler(w, req, id)
			} else if req.Method == http.MethodDelete {
				ws.deleteWalletHandler(w, req, id)
			} else if req.Method == http.MethodPut {
				ws.updateWalletHandler(w, req, id)
			} else {
				http.Error(w, fmt.Sprintf("expect method GET or DELETE at /wallets/{id}, got %v", req.Method), http.StatusMethodNotAllowed)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			//финансовые операции с кошельком
		} else if len(pathParams) == 3 {
			if pathParams[2] == "deposit" {
				ws.depositWalletHandler(w, req, id)
			} else if pathParams[2] == "withdraw" {
				ws.withdrawWalletHandler(w, req, id)
			} else if pathParams[2] == "transfer" {
				ws.transferWalletHandler(w, req, id)
			} else {
				http.Error(w, fmt.Sprintf("expect method with path '/deposit', '/withdraw' or '/transfer'  at /wallets/{id}, got %v", req.Method), http.StatusMethodNotAllowed)
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}
}

// POST /wallets - создание кошелька
func (ws *walletServer) createWalletHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	type RequestWallet struct {
		Name string `json:"name"`
	}

	type ResponseWallet struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	handlingRequest(w, req)

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var rw RequestWallet
	if err := decode.Decode(&rw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wallet, err := ws.store.CreateWallet(rw.Name)
	if err == nil {
		renderJSON(w, ResponseWallet{wallet.Id, wallet.Name, wallet.Status})
		kafka.WalletСreateEvent(wallet.Id, wallet.Status)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	t := time.Since(start)
	log.Printf("handling wallet create at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// GET /wallets - получение списка всех доступных кошельков (в т.ч. деактивированных)
func (ws *walletServer) getAllWalletsHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	allWallets := ws.store.GetAllWallets()
	renderJSON(w, allWallets)

	t := time.Since(start)
	log.Printf("handling get all wallets at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// GET /wallets/{id} - получение кошелька по его идентификатору
func (ws *walletServer) getWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	wallet, err := ws.store.GetWallet(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, wallet)

	t := time.Since(start)
	log.Printf("handling get wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// DELETE /wallet/{id} - деактивация кошелька по идентификатору в path params
func (ws *walletServer) deleteWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	type ResponseWallet struct {
		Success    bool   `json:"success"`
		Error_code string `json:"error_code"`
		Id         string `json:"id"`
	}
	wallet, err := ws.store.DeleteWallet(id)
	if err == nil {
		renderJSON(w, ResponseWallet{Success: true, Error_code: "OK", Id: id})
		kafka.WalletDeleteEvent(wallet.Id, wallet.Status)
	} else {
		w.WriteHeader(http.StatusNotFound)
		renderJSON(w, ResponseWallet{Success: false, Error_code: err.Error(), Id: id})
	}

	t := time.Since(start)
	log.Printf("handling delete wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// PUT /wallets/{id} - обновление кошелька по его идентификатору
func (ws *walletServer) updateWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	type RequestWallet struct {
		Name string `json:"name"`
	}

	type ResponseWallet struct {
		Success    bool   `json:"success"`
		Error_code string `json:"error_code"`
		Id         string `json:"id"`
	}

	handlingRequest(w, req)

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var rw RequestWallet
	if err := decode.Decode(&rw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	
	error := ws.store.UpdateWallet(rw.Name, id)
	if error == nil {
		renderJSON(w, ResponseWallet{Success: true, Error_code: "OK", Id: id})
	} else if error.Error() == "'Name' must consist of more than one symbol"{
		renderJSON(w, ResponseWallet{Success: false, Error_code: error.Error(), Id: id})
		w.WriteHeader(http.StatusBadRequest)
	} else if error.Error() == "Wallet not found"{
		renderJSON(w, ResponseWallet{Success: false, Error_code: error.Error(), Id: id})
		w.WriteHeader(http.StatusNotFound)
	}

	t := time.Since(start)
	log.Printf("handling update wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// POST /wallets/{id}/deposit - метод пополнения кошелька
func (ws *walletServer) depositWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	type RequestWallet struct {
		Amount float64 `json:"amount"`
	}

	type ResponseWallet struct {
		Success    bool    `json:"success"`
		Error_code string  `json:"error_code"`
		Amount     float64 `json:"amount"`
	}

	handlingRequest(w, req)

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var rw RequestWallet
	if err := decode.Decode(&rw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := ws.store.DepositWallet(rw.Amount, id)
	if err == nil {
		renderJSON(w, ResponseWallet{Success: true, Error_code: "OK", Amount: rw.Amount})
		kafka.WalletDepositeEvent(rw.Amount)
	} else if err.Error() == "Wallet  not found" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: err.Error(), Amount: rw.Amount})
		w.WriteHeader(http.StatusNotFound)
	} else if err.Error() == "Amount cannot be negative" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: err.Error(), Amount: rw.Amount})
		w.WriteHeader(http.StatusBadRequest)
	}

	t := time.Since(start)
	log.Printf("handling deposit wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())

}

// POST /wallets/{id}/withdraw - метод снятия средств с кошелька
func (ws *walletServer) withdrawWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	type RequestWallet struct {
		Amount float64 `json:"amount"`
	}

	type ResponseWallet struct {
		Success    bool    `json:"success"`
		Error_code string  `json:"error_code"`
		Amount     float64 `json:"amount"`
	}

	handlingRequest(w, req)

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var rw RequestWallet
	if err := decode.Decode(&rw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := ws.store.WithdrawWallet(rw.Amount, id)
	if err == nil {
		renderJSON(w, ResponseWallet{Success: true, Error_code: "OK", Amount: rw.Amount})
		kafka.WalletWithdrawEvent(rw.Amount)
	} else if err.Error() == "Amount cannot be negative" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: "Amount cannot be negative", Amount: rw.Amount})
		w.WriteHeader(http.StatusBadRequest)
		} else if err.Error() == "Wallet not found" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: "Wallet not found", Amount: rw.Amount})
		w.WriteHeader(http.StatusNotFound)
	} else if err.Error() == "Insufficient funds in the wallet" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: "Insufficient funds in the wallet", Amount: rw.Amount})
		w.WriteHeader(http.StatusBadRequest)
	}

	t := time.Since(start)
	log.Printf("handling withdraw wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

// POST /wallets/{id}/transfer - перевод между двумя кошельками.
func (ws *walletServer) transferWalletHandler(w http.ResponseWriter, req *http.Request, id string) {
	start := time.Now()

	type RequestWallet struct {
		Amount      float64 `json:"amount"`
		Transfer_to string  `json:"transfer_to"`
	}

	type ResponseWallet struct {
		Success    bool    `json:"success"`
		Error_code string  `json:"error_code"`
		Amount     float64 `json:"amount"`
	}

	handlingRequest(w, req)

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var rw RequestWallet
	if err := decode.Decode(&rw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	error := ws.store.TransferWallet(rw.Amount, id, rw.Transfer_to)
	if error == nil {
		renderJSON(w, ResponseWallet{Success: true, Error_code: "OK", Amount: rw.Amount})
		kafka.WalletTransferEvent(rw.Amount)
	} else if error.Error() == "Sender wallet not found" || error.Error() == "Recipient wallet not found"{
		renderJSON(w, ResponseWallet{Success: false, Error_code: error.Error(), Amount: rw.Amount})
		w.WriteHeader(http.StatusNotFound)
	} else if error.Error() == "Insufficient funds in the sender's wallet" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: error.Error(), Amount: rw.Amount})
		w.WriteHeader(http.StatusBadRequest)
	} else if error.Error() == "Amount cannot be negative" {
		renderJSON(w, ResponseWallet{Success: false, Error_code: error.Error(), Amount: rw.Amount})
		w.WriteHeader(http.StatusBadRequest)
	}

	t := time.Since(start)
	log.Printf("handling transfer wallet at %s (request processing time: %d microseconds)\n", req.URL.Path, t.Microseconds())
}

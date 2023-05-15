package main

import (
	"net/http"
	"os"
	"os/signal"
	"log"
	"fmt"
	"time"
	"encoding/json"

	"stats_wallets.com/kafka"
	"stats_wallets.com/stats"

)

type statsServer struct {
	store *stats.StatsStore
}

func NewStatsServer() *statsServer {
	store := stats.New()
	return &statsServer{store: store}
}

func (ss *statsServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/stats/wallets" {
		if r.Method == http.MethodGet{
			ss.getStatsHandler(w, r)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET at /stats/wallets, got %v", r.Method), http.StatusMethodNotAllowed)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
	//возможно нужен доп else
}

func (ss *statsServer) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	//stats.Lock.Lock()
	//defer stats.Lock.Unlock()

	fmt.Printf("GET STATSSSSSS: \n", stats.ResultStats.Wallet.Active)
	response, err := json.Marshal(stats.ResultStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	t := time.Since(start)
	log.Printf("handling get stats of wallets at %s (request processing time: %d microseconds)\n", r.URL.Path, t.Microseconds())

}

func main() {
	mux := http.NewServeMux()
	server := NewStatsServer()
	mux.HandleFunc("/stats/wallets", server.statsHandler)
	
	kafka.StartConsumer()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		log.Println("Received interrupt signal. Shutting down...")
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":5300", mux))
}
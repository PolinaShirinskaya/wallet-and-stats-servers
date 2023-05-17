package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"stats_wallets.com/kafka"
	"stats_wallets.com/stats"
)


func main() {
	
	http.HandleFunc("/wallets/stats/", statsHandler)

	go func() {
		kafka.StartConsumer()
	}()

	log.Fatal(http.ListenAndServe(":5200", nil))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	getWalletStats(w, r)
}

func getWalletStats(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
	
		stats.Lock.Lock()
		defer stats.Lock.Unlock()
	
		response, err := json.Marshal(stats.ResultStats)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	
		t := time.Since(start)
		log.Printf("handling get stats of wallets at %s (request processing time: %d microseconds)\n", r.URL.Path, t.Microseconds())
	
}

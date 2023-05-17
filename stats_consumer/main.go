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

























// package main

// import (
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"log"
// 	//"fmt"
// 	"time"
// 	"encoding/json"

// 	"stats_wallets.com/kafka"
// 	"stats_wallets.com/stats"

// )

// // type statsServer struct {
// // 	store *stats.StatsStore
// // }

// // func NewStatsServer() *statsServer {
// // 	store := stats.New()
// // 	return &statsServer{store: store}
// // }

// // func (ss *statsServer) statsHandler(w http.ResponseWriter, r *http.Request) {
// // 	fmt.Println("STATS_HANDLER!!!!")
// // 	if r.URL.Path == "/wallets/stats" {
// // 		if r.Method == http.MethodGet{
// // 			ss.getStatsHandler(w, r)
// // 		} else {
// // 			http.Error(w, fmt.Sprintf("expect method GET at /wallets/stats, got %v", r.Method), http.StatusMethodNotAllowed)
// // 			w.WriteHeader(http.StatusMethodNotAllowed)
// // 		}
// // 	} else {
// // 		fmt.Println("ERROR:::STATS_HANDLER!!!!")
// // 	}
// // }

// //func (ss *statsServer) getStatsHandler(w http.ResponseWriter, r *http.Request) {
// func getWalletStats(w http.ResponseWriter, r *http.Request) {
	
// 	start := time.Now()

// 	stats.Lock.Lock()
// 	defer stats.Lock.Unlock()

// 	response, err := json.Marshal(stats.ResultStats)
// 	if err != nil {
// 		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
// 		return
// 	}

// 	// fmt.Print("GET STATSSSSSS: \n", stats.ResultStats.Wallet.Active)
// 	// response, err := json.Marshal(stats.ResultStats)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(response)

// 	t := time.Since(start)
// 	log.Printf("handling get stats of wallets at %s (request processing time: %d microseconds)\n", r.URL.Path, t.Microseconds())

// }

// func main() {
// 	mux := http.NewServeMux()
// 	//server := NewWalletServer()
// 	mux.HandleFunc("/wallets/stats", getWalletStats)

// 	kafka.StartConsumer()

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, os.Interrupt)

// 	go func() {
// 		<-signals
// 		log.Println("Received interrupt signal. Shutting down...")
// 		os.Exit(0)
// 	}()

// 	log.Fatal(http.ListenAndServe(":5200", mux))
	
// 	/*mux := http.NewServeMux()
// 	// server := NewStatsServer()
// 	//mux.HandleFunc("/wallets/stats", server.statsHandler)
	
// 	mux.HandleFunc("/wallets/stats/", getWalletStats)

// 	kafka.StartConsumer()

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, os.Interrupt)

// 	go func() {
// 		<-signals
// 		log.Println("Received interrupt signal. Shutting down...")
// 		os.Exit(0)
// 	}()


// 	err := http.ListenAndServe(":5200", mux)
// 	if err != nil {
// 		log.Fatal("Server error:", err)
// 	}
// 	//log.Fatal(http.ListenAndServe(":5200", mux))*/
// }
package stats

import "sync"

type StatsWallets struct {
	Id 		string
	Status 	string
}

type StatsStore struct {
	sync.Mutex

	stats []StatsWallets
}

func New() *StatsStore {
	ss := &StatsStore{}
	ss.stats = make([]StatsWallets, 0)
	return ss
}
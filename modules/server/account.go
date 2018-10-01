package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/zc-staff/openstock/api"
)

var (
	users = make(map[uint64]string)
	vis   = map[string]bool{ipo: true}
	lock  sync.RWMutex
)

func randint64() uint64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint64(b[:])
}

type registerInput struct {
	Name string
}

type bidInput struct {
	Token         uint64
	Price, Amount int
}

type queryInput struct {
	Token uint64
}

func fromJSON(w http.ResponseWriter, r *http.Request, obj interface{}) bool {
	dec := json.NewDecoder(r.Body)
	if dec.Decode(obj) != nil {
		http.Error(w, "argument error", 400)
		return false
	}
	return true
}

func register(w http.ResponseWriter, r *http.Request) {
	var inp registerInput
	if !fromJSON(w, r, &inp) {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	if vis[inp.Name] {
		http.Error(w, "duplicate account", 400)
		return
	}
	token := randint64()
	vis[inp.Name] = true
	users[token] = inp.Name

	// NOTE: there should not be error
	bank.AddAccount(inp.Name, initial)
	company.AddAccount(inp.Name, 0)

	toJSON(w, map[string]uint64{"token": token})
}

func queryAccount(w http.ResponseWriter, r *http.Request) {
	var inp queryInput
	if !fromJSON(w, r, &inp) {
		return
	}

	lock.RLock()
	defer lock.RUnlock()

	if name, ok := users[inp.Token]; ok {
		// NOTE: there should not be error
		bf, ba, _ := bank.Query(name)
		cf, ca, _ := company.Query(name)
		toJSON(w, map[string]map[string]int{
			"bank":    {"available": ba, "frozen": bf},
			"company": {"available": ca, "frozen": cf},
		})
		return
	}
	http.Error(w, "account not exist", 400)
}

func bidHandler(sell bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var inp bidInput
		if !fromJSON(w, r, &inp) {
			return
		}

		lock.RLock()
		defer lock.RUnlock()

		if name, ok := users[inp.Token]; ok {
			app.Post(api.OnBid, api.Bid{
				Bidder: name, Sell: sell,
				Price: inp.Price, Amount: inp.Amount,
			})

			return
		}
		http.Error(w, "account not exist", 400)
	}
}

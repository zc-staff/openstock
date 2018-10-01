package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zc-staff/openstock/api"
)

var (
	market        api.Market
	bank, company api.Saving
	periods       int
	nextPeriod    time.Time
)

func toJSON(w http.ResponseWriter, obj interface{}) {
	enc := json.NewEncoder(w)
	err := enc.Encode(obj)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func transaction(w http.ResponseWriter, r *http.Request) {
	toJSON(w, market.GetTransactions())
}

func bid(w http.ResponseWriter, r *http.Request) {
	buys, sells := market.GetBids()
	ret := map[string][]api.Bid{"buys": buys, "sells": sells}
	toJSON(w, ret)
}

func queryTime(w http.ResponseWriter, r *http.Request) {
	toJSON(w, map[string]interface{}{
		"periods": periods, "nextPeriod": nextPeriod, "now": time.Now(),
	})
}

package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
)

var (
	interest float64
)

func onPeriod(interface{}) {
	lock.Lock()
	defer lock.Unlock()
	// NOTE: there should not be error

	// interest
	interest = smoothNoise()
	log.Println("new interest", interest)
	for _, name := range users {
		_, m, _ := bank.Query(name)
		_, s, _ := company.Query(name)

		inc := int(math.Floor(float64(m) * bankRate))
		inc += int(math.Floor(float64(s*value) * interest))
		bank.Income(name, inc)
	}

	// clear bids
	log.Println("clear bids")
	market.Clear()

	// unfreeze
	log.Println("unfreeze accounts")
	for _, name := range users {
		bank.Unfreeze(name)
		company.Unfreeze(name)
	}
}

func onIPO(interface{}) {
	market.Sell(ipo, 0, stocks)
}

func whiteNoise() float64 {
	return math.Max(0, rand.NormFloat64()*compStd+compRate)
}

func smoothNoise() float64 {
	mean := 0.5 * (compRate + interest)
	std := compStd * math.Sqrt(0.75)
	return math.Max(0, rand.NormFloat64()*std+mean)
}

func queryInterest(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	toJSON(w, interest)
}

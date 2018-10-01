package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/zc-staff/openplayer/api"
	stock "github.com/zc-staff/openstock/api"
)

var (
	app api.Application
)

type entry struct{}

func (e *entry) Load(ap api.Application) error {
	app = ap
	app.Subscribe(api.OnInit, func(interface{}) {
		// modules
		market = api.EnsureModule(app, "market").(stock.Market)
		bank = api.EnsureModule(app, "bank").(stock.Saving)
		company = api.EnsureModule(app, "company").(stock.Saving)

		// ipo account
		bank.AddAccount(ipo, 0)
		company.AddAccount(ipo, stocks)

		// period
		seed := int64(randint64())
		rand.Seed(seed)
		log.Println("seed", seed)
		interest = whiteNoise()
		log.Println("initial interest", interest)
		go func() {
			tick := time.Tick(period)
			for {
				nextPeriod = time.Now().Add(period)
				<-tick
				if periods == ipoTime {
					app.Post(stock.OnIPO, nil)
				}
				app.Post(stock.OnPeriod, nil)
				periods++
			}
		}()

		// start server
		http.HandleFunc("/api/transaction", transaction)
		http.HandleFunc("/api/bid", bid)
		http.HandleFunc("/api/register", register)
		http.HandleFunc("/api/account", queryAccount)
		http.HandleFunc("/api/interest", queryInterest)
		http.HandleFunc("/api/time", queryTime)
		http.HandleFunc("/api/buy", bidHandler(false))
		http.HandleFunc("/api/sell", bidHandler(true))
		go func() {
			log.Fatal(http.ListenAndServe(":8887", nil))
		}()
	})
	app.Subscribe(stock.OnPeriod, onPeriod)
	app.Subscribe(stock.OnIPO, onIPO)
	return nil
}

var Entry entry

func main() {}

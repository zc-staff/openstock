package main

import (
	"log"

	"github.com/zc-staff/openplayer/api"
	stock "github.com/zc-staff/openstock/api"
)

type entry struct{}

func (e *entry) Load(app api.Application) error {
	mark.start()
	app.RegisterModule("market", &mark)
	app.Subscribe(api.OnInit, func(interface{}) {
		bank = api.EnsureModule(app, "bank").(stock.Saving)
		company = api.EnsureModule(app, "company").(stock.Saving)
	})
	app.Subscribe(stock.OnBid, func(arg interface{}) {
		bid := arg.(stock.Bid)
		var err error
		if bid.Sell {
			err = mark.Sell(bid.Bidder, bid.Price, bid.Amount)
		} else {
			err = mark.Buy(bid.Bidder, bid.Price, bid.Amount)
		}
		if err != nil {
			log.Println(bid.Bidder, "bid failed:", err)
		}
	})
	return nil
}

var Entry entry

func main() {}

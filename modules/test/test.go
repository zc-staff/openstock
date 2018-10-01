package main

import (
	"log"

	"github.com/zc-staff/openplayer/api"
	stock "github.com/zc-staff/openstock/api"
)

type entry struct {
	app api.Application
}

func (t *entry) Load(app api.Application) error {
	t.app = app
	app.Subscribe(api.OnInit, t.run)
	return nil
}

func printt(m stock.Market, b stock.Saving, c stock.Saving) {
	buys, sells := m.GetBids()
	log.Println("buys", buys)
	log.Println("sells", sells)
	log.Println("trans", m.GetTransactions())
	log.Println(b.Query("tiny"))
	log.Println(c.Query("tiny"))
}

func (t *entry) run(interface{}) {
	log.Println("hello")

	m, err := t.app.GetModule("market")
	if err != nil {
		log.Fatal(err)
	}
	market := m.(stock.Market)

	b, err := t.app.GetModule("bank")
	if err != nil {
		log.Fatal(err)
	}
	c, err := t.app.GetModule("company")
	if err != nil {
		log.Fatal(err)
	}
	bank, company := b.(stock.Saving), c.(stock.Saving)

	bank.AddAccount("tiny", 1000)
	bank.AddAccount("zzy", 0)
	company.AddAccount("tiny", 0)
	company.AddAccount("zzy", 1000)
	printt(market, bank, company)

	log.Println("tiny buy 100 @ 2", market.Buy("tiny", 2, 100))
	printt(market, bank, company)
	log.Println("zzy sell 100 @ 0", market.Sell("zzy", 0, 100))
	printt(market, bank, company)
	log.Println("zzy sell 100 @ 2", market.Sell("zzy", 2, 100))
	printt(market, bank, company)
	log.Println("tiny buy 100 @ 3", market.Buy("tiny", 8, 100))
	printt(market, bank, company)
}

var Entry entry

func main() {}

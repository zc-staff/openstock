package main

import (
	"container/heap"
	"log"
	"time"

	"github.com/zc-staff/openstock/api"
)

var (
	bank, company api.Saving
	mark          market
)

type market struct {
	buys, sells bidList
	trans       []api.Transaction
}

func (m *market) start() {
	m.trans = []api.Transaction{}
	m.buys, m.sells = bidList{}, bidList{}
}

func (m *market) deal(bid api.Bid) error {
	if bid.Price < 0 || bid.Amount <= 0 {
		return api.ArgError
	}

	var err error
	if bid.Sell {
		err = company.Freeze(bid.Bidder, bid.Amount)
	} else {
		err = bank.Freeze(bid.Bidder, bid.Amount*bid.Price)
	}
	if err != nil {
		return err
	}

	me, other := &m.buys, &m.sells
	if bid.Sell {
		me, other = &m.sells, &m.buys
	}

	ret, am := other.clear(bid)
	if am > 0 {
		bid.Amount = am
		heap.Push(me, bid)
	}

	for _, r := range ret {
		// NOTE: there should not be error
		if bid.Sell {
			company.Deduct(bid.Bidder, r.Amount)
			bank.Income(bid.Bidder, r.Amount*r.Price)
			bank.Deduct(r.Bidder, r.Amount*r.Price)
			company.Income(r.Bidder, r.Amount)
		} else {
			bank.Deduct(bid.Bidder, r.Amount*r.Price)
			company.Income(bid.Bidder, r.Amount)
			company.Deduct(r.Bidder, r.Amount)
			bank.Income(r.Bidder, r.Amount*r.Price)
		}

		t := api.Transaction{
			Buyer: bid.Bidder, Seller: r.Bidder, Time: time.Now(),
			Price: r.Price, Amount: r.Amount,
		}
		if bid.Sell {
			t.Buyer, t.Seller = t.Seller, t.Buyer
		}
		m.trans = append(m.trans, t)
	}
	return nil
}

func (m *market) Clear() error {
	// NOTE: there should not be error
	m.buys = m.buys[:0]
	m.sells = m.sells[:0]
	return nil
}

func (m *market) Buy(name string, price int, amount int) error {
	log.Println(name, "buy", amount, "@", price)
	return m.deal(api.Bid{
		Bidder: name, Time: time.Now(),
		Price: price, Amount: amount,
	})
}

func (m *market) Sell(name string, price int, amount int) error {
	log.Println(name, "sell", amount, "@", price)
	return m.deal(api.Bid{
		Bidder: name, Time: time.Now(), Sell: true,
		Price: price, Amount: amount,
	})
}

func (m *market) GetTransactions() []api.Transaction {
	return m.trans
}

func (m *market) GetBids() ([]api.Bid, []api.Bid) {
	return m.buys, m.sells
}

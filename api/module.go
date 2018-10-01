package api

import (
	"errors"
	"time"
)

var (
	ArgError         = errors.New("argument error")
	AccountNotExist  = errors.New("account not exist")
	AccountExist     = errors.New("account exist")
	InsufficientFund = errors.New("insufficient fund")
)

type Transaction struct {
	Buyer, Seller string
	Price, Amount int
	Time          time.Time
}

type Bid struct {
	Bidder        string
	Price, Amount int
	Sell          bool
	Time          time.Time
}

type Market interface {
	GetTransactions() []Transaction
	GetBids() ([]Bid, []Bid)
	Buy(string, int, int) error
	Sell(string, int, int) error
	Clear() error
}

type Saving interface {
	AddAccount(string, int) error
	Query(string) (int, int, error)
	Freeze(string, int) error
	Unfreeze(string) error
	Deduct(string, int) error
	Income(string, int) error
}

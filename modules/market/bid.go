package main

import (
	"container/heap"

	"github.com/zc-staff/openstock/api"
)

type bidList []api.Bid

func (b bidList) Len() int {
	return len(b)
}

func (b bidList) Less(i, j int) bool {
	if b[i].Price == b[j].Price {
		return b[i].Time.Before(b[j].Time)
	}
	return b[i].Sell == (b[i].Price < b[j].Price)
}

func (b bidList) Swap(i, j int) {
	t := b[i]
	b[i] = b[j]
	b[j] = t
}

func (b *bidList) Push(x interface{}) {
	*b = append(*b, x.(api.Bid))
}

func (b *bidList) Pop() interface{} {
	old := *b
	n := len(old)
	*b = old[:n-1]
	return old[n-1]
}

func (b *bidList) clear(bid api.Bid) ([]api.Bid, int) {
	var ret []api.Bid
	for len(*b) > 0 {
		now := &(*b)[0]
		if (bid.Sell && bid.Price > now.Price) || (!bid.Sell && bid.Price < now.Price) {
			break
		}
		if now.Amount >= bid.Amount {
			bb := *now
			bb.Amount = bid.Amount
			ret = append(ret, bb)
			now.Amount -= bid.Amount
			if now.Amount <= 0 {
				heap.Pop(b)
			}
			return ret, 0
		}
		bid.Amount -= now.Amount
		ret = append(ret, *now)
		heap.Pop(b)
	}
	return ret, bid.Amount
}

package main

import (
	"log"

	"github.com/zc-staff/openstock/api"
)

type saving struct {
	name          string
	avail, frozen map[string]int
}

var (
	bank    = saving{name: "bank"}
	company = saving{name: "company"}
)

func (s *saving) start() {
	s.avail = make(map[string]int)
	s.frozen = make(map[string]int)
}

func (s *saving) AddAccount(name string, num int) error {
	log.Println(s.name, "add account", name, num)
	if _, ok := s.avail[name]; ok {
		return api.AccountExist
	}
	s.avail[name], s.frozen[name] = num, 0
	return nil
}

func (s *saving) Query(name string) (int, int, error) {
	log.Println(s.name, "query", name)
	if a, ok := s.avail[name]; ok {
		return s.frozen[name], a, nil
	}
	return 0, 0, api.AccountNotExist
}

func (s *saving) Freeze(name string, num int) error {
	log.Println(s.name, "freeze", name, num)
	if _, ok := s.avail[name]; ok {
		if s.avail[name] < num {
			return api.InsufficientFund
		}
		s.avail[name] -= num
		s.frozen[name] += num
		return nil
	}
	return api.AccountNotExist
}

func (s *saving) Unfreeze(name string) error {
	log.Println(s.name, "unfreeze", name)
	if _, ok := s.avail[name]; ok {
		s.avail[name] += s.frozen[name]
		s.frozen[name] = 0
		return nil
	}
	return api.AccountNotExist
}

func (s *saving) Deduct(name string, num int) error {
	log.Println(s.name, "deduct", name, num)
	if _, ok := s.avail[name]; ok {
		if s.frozen[name] < num {
			return api.InsufficientFund
		}
		s.frozen[name] -= num
		return nil
	}
	return api.AccountNotExist
}

func (s *saving) Income(name string, num int) error {
	log.Println(s.name, "income", name, num)
	if _, ok := s.avail[name]; ok {
		s.frozen[name] += num
		return nil
	}
	return api.AccountNotExist
}

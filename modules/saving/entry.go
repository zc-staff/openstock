package main

import "github.com/zc-staff/openplayer/api"

type entry struct{}

func (e *entry) Load(app api.Application) error {
	bank.start()
	company.start()
	app.RegisterModule("bank", &bank)
	app.RegisterModule("company", &company)
	return nil
}

var Entry entry

func main() {}

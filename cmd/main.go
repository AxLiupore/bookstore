package main

import (
	"bookstore/conf"
	"bookstore/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	err := r.Run(conf.HttpPort)
	if err != nil {
		return
	}
}

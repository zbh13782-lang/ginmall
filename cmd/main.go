package main

import (
	"gin_mall/conf"
	"gin_mall/loading"
	"gin_mall/routes"
)

func main() {
	conf.Init()
	loading.Loading()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}

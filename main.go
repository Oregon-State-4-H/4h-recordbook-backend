package main

import (
	. "4h-recordbook/backend/config"
	"4h-recordbook/backend/internal/router"
)

var config Config

func init() {
	config.Read()
}

// @title	4H Record Books API
func main(){

	rtr := router.New()
  	rtr.Run(config.Server.Port)
	
}
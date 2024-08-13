package main

import (
	"ecommerce-golang/config"
	"ecommerce-golang/internal/api"
	"log"
)

func main() {

	cfg , err:= config.SetupEnv()

	if err != nil {
		log.Fatalf("config file is not loaded %v", err)
	}

	api.StartServer(cfg)
}

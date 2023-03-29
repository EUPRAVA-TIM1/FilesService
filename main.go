package main

import (
	"FileService/config"
	"FileService/startup"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(&config)
	server.Start()
}

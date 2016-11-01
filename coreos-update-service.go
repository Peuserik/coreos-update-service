package main

import "github.com/pegerto/coreos-update-service/coreos"

func main() {
	service := coreos.NewServer()
	service.Server()
}



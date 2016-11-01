package main

import "./coreos"

func main() {
	service := coreos.NewServer()
	service.Server()
}

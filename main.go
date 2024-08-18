package main

import myNetwork "go-chat/network"

func main() {
	network := myNetwork.NewServer()
	network.StartServer()
}

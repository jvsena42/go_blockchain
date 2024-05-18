package main

import (
	"fmt"

	"github.com/jvsena42/go_blockchain/utils"
)

func main() {
	//fmt.Println(utils.IsFoundNode("127.0.0.1", 3333))
	//fmt.Println(utils.IsFoundNode("localhost", 3333))
	//fmt.Println(utils.FindNeighbors("127.0.0.1", 3333, 0, 3, 3333, 3336))

	myAddress := utils.GetHost()
	fmt.Println(utils.FindNeighbors(myAddress, 3333, 0, 3, 3333, 3336))
}

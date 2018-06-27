package main

import (
	"fmt"
	"LOG735-PG/src/app"
)

func main() {
	fmt.Println("Hello")

	miner := app.Miner{}
	miner.CreateBlock()
	fmt.Println(miner.Chain)
	
}
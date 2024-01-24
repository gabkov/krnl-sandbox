package client

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClient() (*ethclient.Client) {
	client, err := ethclient.Dial("http://127.0.0.1:8545") // hardhat local node
	if err != nil {
		log.Fatal(err)
	}

	return client
}
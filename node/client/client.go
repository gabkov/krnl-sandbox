package client

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
Under the hood we are connecting to a local hardhat node 
and extending it with krnl specific rpc calls
*/
func GetClient() (*ethclient.Client) {
	client, err := ethclient.Dial("http://127.0.0.1:8545") // hardhat local node
	if err != nil {
		log.Fatal(err)
	}

	return client
}
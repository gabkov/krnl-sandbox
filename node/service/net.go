package service

import (
	"context"
	"log"
	"math/big"

	"github.com/gabkov/krnl-node/client"
)

type Net struct{}

func (n *Net) Version() (*big.Int, error) {
	log.Println("net_version")
	
	client := client.GetClient()

	return client.NetworkID(context.Background())
}
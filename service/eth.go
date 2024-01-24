package service

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type Eth struct{}

func (t *Eth) ChainId() int {
	return 31337 // hardhat network
}

func (t *Eth) GetBlockByNumber(blockTag string, includeTx bool) (map[string]interface{}, error) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	var head *types.Header
	err = client.Client().CallContext(context.Background(), &head, "eth_getBlockByNumber", blockTag, includeTx)
	if err != nil {
		log.Println("can't get latest block:", err)
		return make(map[string]interface{}), err
	}

	block, err := client.BlockByNumber(context.Background(), nil) // nil means lates and ethers.js asks for latest
	if err != nil {
		log.Println(err)
	}

	return RPCMarshalBlock(block.WithSeal(head), true, true, &params.ChainConfig{}), nil
}

func (t *Eth) GetTransactionCount(account common.Address, blockTag string) (uint64, error) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	return client.NonceAt(context.Background(), account, nil)
}

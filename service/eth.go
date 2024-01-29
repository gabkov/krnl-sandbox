package service

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gabkov/krnl-node/client"
)

type Eth struct{}

func (t *Eth) ChainId() int {
	log.Println("eth_chainId called")
	return 31337 // hardhat network
}

func (t *Eth) GetBlockByNumber(blockTag string, includeTx bool) (map[string]interface{}, error) {
	log.Println("eth_getBlockByNumber called")
	client := client.GetClient()

	var head *types.Header
	err := client.Client().CallContext(context.Background(), &head, "eth_getBlockByNumber", blockTag, includeTx)
	if err != nil {
		log.Println("can't get latest block:", err)
		return make(map[string]interface{}), err
	}

	block, err := client.BlockByNumber(context.Background(), nil) // nil means latest and ethers.js asks for latest
	if err != nil {
		log.Println(err)
	}

	return RPCMarshalBlock(block.WithSeal(head), true, true, &params.ChainConfig{}), nil
}

func (t *Eth) GetTransactionCount(account common.Address, blockTag string) (uint64, error) {
	log.Println("eth_getTransactionCount called")
	client := client.GetClient()

	return client.NonceAt(context.Background(), account, nil) // nil means latest and ethers.js asks for latest
}

func (t *Eth) EstimateGas(ethCallMsg map[string]interface{}) (uint64, error) {
	log.Println("eth_estimateGas called")
	client := client.GetClient()

	var hex hexutil.Uint64
	err := client.Client().CallContext(context.Background(), &hex, "eth_estimateGas", ethCallMsg)

	if err != nil {
		return 0, err
	}

	return uint64(hex), nil
}

func (t *Eth) Call(ethCallMsg map[string]interface{}, blockTag string) (string, error) {
	log.Println("eth_call called")
	client := client.GetClient()

	var hex hexutil.Bytes
	err := client.Client().CallContext(context.Background(), &hex, "eth_call", ethCallMsg, nil)
	if err != nil {
		return "", err
	}

	return hex.String(), nil
}

func (t *Eth) SendRawTransaction(tx string) (string, error) {
	log.Println("eth_sendRawTransaction called")
	client := client.GetClient()

	rawTxBytes, err := hex.DecodeString(tx[2:])

	txparsed := new(types.Transaction)


	err = txparsed.UnmarshalBinary(rawTxBytes)
    if err != nil {
        log.Println("err:", err)
    }

	err = client.SendTransaction(context.Background(), txparsed)
	if err != nil {
		log.Println(err)
		return "" ,err
	}

	return txparsed.Hash().Hex(), nil
}

func (t *Eth) BlockNumber() (uint64, error) {
	log.Println("eth_blockNumber called")
	client := client.GetClient()
	return client.BlockNumber(context.Background())
}

func (t *Eth) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	log.Println("eth_getTransactionReceipt called")
	client := client.GetClient()
	return client.TransactionReceipt(context.Background(), txHash)
}

package service

import (
	"context"
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gabkov/krnl-node/client"
	"github.com/gabkov/krnl-node/faas"
)

type Eth struct{}

/*
Pauses the transaction and check's if there was any additional 
data (FaaS requests) concatenated to the end of the input data field.
*/
func (t *Eth) SendRawTransaction(rawTx string) (string, error) {
	log.Println("eth_sendRawTransaction")

	client := client.GetClient()

	rawTxBytes, _ := hex.DecodeString(rawTx[2:])

	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(rawTxBytes)
	if err != nil {
		log.Fatal("err:", err)
	}

	// simulate stopping tx here
	// grabbing the requested FaaS services from the end of the input-data
	separator := "000000000000000000000000000000000000000000000000000000000000003a" // :
	res := strings.Split(hexutil.Encode(tx.Data()), separator)

	// if len is more than 1 some message is concatenated to the end of the input-data
	if len(res) > 1 {
		for i := 1; i < len(res); i++ {
			faasRequest, err := hex.DecodeString(res[i])
			if err != nil {
				log.Fatal(err)
			}
			// mock FaaS service call
			err = faas.CallService(faasRequest, tx)
			if err != nil {
				log.Println("FaaS request failed:", err)
				return "", err
			}
		}
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("tx sent: %s", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

func (t *Eth) ChainId() (string)  {
	// log.Println("eth_chainId") // commenting it out too many logs
	
	return "0x7a69"
}

func (t *Eth) GasPrice() (*big.Int, error) {
	log.Println("eth_gasPrice")
	client := client.GetClient()

	return client.SuggestGasPrice(context.Background())
}

func (t *Eth) GetBalance(account common.Address, blockTag interface{}) (*hexutil.Big, error) {
	log.Println("eth_getBalance")
	client := client.GetClient()

	var _blocktag string
	switch v := blockTag.(type) {
		case string:
			_blocktag = v
		case float64:
			_blocktag = toBlockNumArg(big.NewInt(int64(v)))
	}

	var result hexutil.Big
	err := client.Client().CallContext(context.Background(), &result, "eth_getBalance", account, _blocktag)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	return &result, nil
}

func (t *Eth) GetBlockByNumber(blockTag interface{}, includeTx bool) (map[string]interface{}, error) {
	log.Println("eth_getBlockByNumber")
	
	client := client.GetClient()

	var _blocktag string
	var _block *types.Block
	switch v := blockTag.(type) {
		case string:
			log.Println("string blocktag:",v)
			_blocktag = v
			block, err := client.BlockByNumber(context.Background(), nil) // nil means latest
			if err != nil {
				log.Println("can't get block by number:", err)		
				return make(map[string]interface{}), err
			}
			_block = block
		case float64:
			log.Println("number blocktag:",v)
			_blocktag = toBlockNumArg(big.NewInt(int64(v)))
			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(v)))
			if err != nil {
				log.Println("can't get block by number:", err)		
				return make(map[string]interface{}), err
			}
			_block = block
	}

	var head *types.Header
	err := client.Client().CallContext(context.Background(), &head, "eth_getBlockByNumber", _blocktag, includeTx)
	if err != nil {
		log.Println("can't get block by number:", err)
		return make(map[string]interface{}), err
	}

	return RPCMarshalBlock(_block.WithSeal(head), true, true, &params.ChainConfig{}), nil
}

func (t *Eth) GetBlockByHash(hash common.Hash, includeTx bool) (map[string]interface{}, error) {
	log.Println("eth_getBlockByHash")
	
	client := client.GetClient()

	var head *types.Header
	err := client.Client().CallContext(context.Background(), &head, "eth_getBlockByHash", hash, includeTx)
	if err != nil {
		log.Println("can't get latest block:", err)
		return make(map[string]interface{}), err
	}

	block, err := client.BlockByHash(context.Background(), hash) // nil means latest and ethers.js asks for latest
	if err != nil {
		log.Println(err)
	}

	return RPCMarshalBlock(block.WithSeal(head), true, true, &params.ChainConfig{}), nil
}

func (t *Eth) TransactionByHash(hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	log.Println("eth_getTransactionByHash")
	var json *rpcTransaction

	client := client.GetClient()

	err = client.Client().CallContext(context.Background(), &json, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if json == nil {
		return nil, false, ethereum.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, false, errors.New("server returned transaction without signature")
	}
	// TODO
	// if json.From != nil && json.BlockHash != nil {
	// 	setSenderFromServer(json.tx, *json.From, *json.BlockHash)
	// }
	return json.tx, json.BlockNumber == nil, nil
}

func (t *Eth) GetTransactionCount(account common.Address, blockTag string) (uint64, error) {
	log.Println("eth_getTransactionCount")
	client := client.GetClient()

	var result hexutil.Uint64
	err := client.Client().CallContext(context.Background(), &result, "eth_getTransactionCount", account, blockTag)
	return uint64(result), err
}

func (t *Eth) EstimateGas(ethCallMsg map[string]interface{}) (uint64, error) {
	log.Println("eth_estimateGas")
	client := client.GetClient()

	var hex hexutil.Uint64
	err := client.Client().CallContext(context.Background(), &hex, "eth_estimateGas", ethCallMsg)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return uint64(hex), nil
}

func (t *Eth) FeeHistory(blockCount string, lastBlock string, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	log.Println("eth_feeHistory")

	client := client.GetClient()

	var res feeHistoryResultMarshaling
	if err := client.Client().CallContext(context.Background(), &res, "eth_feeHistory", blockCount, lastBlock, rewardPercentiles); err != nil {
		return nil, err
	}

	reward := make([][]*big.Int, len(res.Reward))
	for i, r := range res.Reward {
		reward[i] = make([]*big.Int, len(r))
		for j, r := range r {
			reward[i][j] = (*big.Int)(r)
		}
	}
	baseFee := make([]*big.Int, len(res.BaseFee))
	for i, b := range res.BaseFee {
		baseFee[i] = (*big.Int)(b)
	}

	return &ethereum.FeeHistory{
		OldestBlock:  (*big.Int)(res.OldestBlock),
		Reward:       reward,
		BaseFee:      baseFee,
		GasUsedRatio: res.GasUsedRatio,
	}, nil
}

func (t *Eth) Call(ethCallMsg map[string]interface{}, blockTag interface{}) (string, error) {
	log.Println("eth_call")

	var _blocktag string
	switch v := blockTag.(type) {
		case string:
			_blocktag = v
		case float64:
			_blocktag = toBlockNumArg(big.NewInt(int64(v)))
	}

	client := client.GetClient()

	var hex hexutil.Bytes
	err := client.Client().CallContext(context.Background(), &hex, "eth_call", ethCallMsg, _blocktag)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return hex.String(), nil
}

func (t *Eth) BlockNumber() (uint64, error) {
	log.Println("eth_blockNumber")
	client := client.GetClient()
	return client.BlockNumber(context.Background())
}

func (t *Eth) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	log.Println("eth_getTransactionReceipt")
	client := client.GetClient()
	return client.TransactionReceipt(context.Background(), txHash)
}

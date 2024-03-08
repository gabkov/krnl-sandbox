package faas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gabkov/krnl-node/build/contracts/krnldapp"
	"github.com/gabkov/krnl-node/client"
)

type TxSendUserOpResponse struct {
	Result string `json:"result" binding:"required"`
}

type TxSendGetUserOpRequest struct {
	JsonRpc string   `json:"jsonrpc"`
	Id      uint16   `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type TxSendGetUserOpResponse struct {
	Params []any `json:"params"`
}

type TxSendSponsorUserOpResult struct {
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

type TxSendSponsorUserOpResponse struct {
	Result TxSendSponsorUserOpResult `json:"result"`
}

// simulate KYT database
// the first address from the local hardhat node config
var kytAddresses = []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}

func CallService(faas string, tx *types.Transaction) error {
	txType := strings.TrimSpace(faas)

	if txType == "KYC" {
		return kyc(tx)
	} else if txType == "KYT" {
		return kyt(tx)
	} else if strings.HasPrefix(txType, "KYT_AA") {
		return kytAA(txType, tx)
	} else {
		log.Println("Unknown function name: ", faas)
		return nil
	}
}

func kyt(tx *types.Transaction) error {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}
	sender := from.Hex()
	for _, value := range kytAddresses {
		if value == sender {
			log.Println("KYT success for address: ", sender)
			return nil
		}
	}

	return errors.New("KYT failed for address " + from.Hex())
}

func kyc(tx *types.Transaction) error {
	log.Println("KYC FaaS not implemented")
	return nil
}

func kytAA(txType string, tx *types.Transaction) error {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}
	sender := from.Hex()

	// get AA request from tx data
	// KYT_AA|hex(request_body)
	separator := "|"
	resBodyBytes, err := hexutil.Decode(strings.Split(txType, separator)[1])

	resBody := string(resBodyBytes)
	log.Printf("Sender %s send %s\n", sender, resBody)

	// receive the OpHash from the bundler
	txSendUserOpResponseBytes, err := callAABundler([]byte(`
	{
	    "jsonrpc": "2.0",
	    "id": 1,
	    "method": "eth_sendUserOperation",
	    "params": [
	        {
	            "sender": "0xd9d567CE0C1BD422424ff8194d7B8D2C6088D452",
	            "nonce": "0x1a",
	            "initCode": "0x",
	            "callData":"0xb61d27f600000000000000000000000020ed044884d83787368861c4f987d9ed7e8aa8a100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002420da7616000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000",
	            "callGasLimit": "0x3076",
	            "verificationGasLimit": "0xee69",
	            "preVerificationGas": "0xc113",
	            "maxFeePerGas": "0x66bb2d8a",
	            "maxPriorityFeePerGas": "0x435a6e80",
	            "paymasterAndData": "0x",
	            "signature": "0x945a0a49468ef4003000def2d305de2bc6d008800c0a0bc5cd2348db998b10d02e9ab6bfca255e22a073c24661f63ee6641efdbdce7bb5c8a274430e5e4dd33f1c"
	        },
	        "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"
	    ]
	}
	`), os.Getenv("ACCOUNT_ABSTRACTION_BUNDLER_ENDPOINT"))
	if err != nil {
		return err
	}

	txSendUserOpResponse := TxSendUserOpResponse{}
	if err := json.Unmarshal(txSendUserOpResponseBytes, &txSendUserOpResponse); err != nil {
		return errors.New("Request to AA Bundler failed")
	}
	// listen for the event from KrnlContract, here is GetCounter
	// if yes, go through
	sepoliaClient := client.GetWsClient()

	// load Krnl contract
	krnlAddr := common.HexToAddress("0x20Ed044884D83787368861C4F987D9ed7e8Aa8A1")
	krnlContract, err := krnldapp.NewKrnldapp(krnlAddr, sepoliaClient)
	if err != nil {
		log.Fatal(err)
	}
	krnlContractAbi, err := abi.JSON(strings.NewReader(string(krnldapp.KrnldappABI)))
	if err != nil {
		log.Fatal(err)
	}

	_ = krnlContract

	query := ethereum.FilterQuery{
		Addresses: []common.Address{krnlAddr},
	}

	logs := make(chan types.Log)

	sub, err := sepoliaClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	// checking log for 30s
	for start := time.Now().Add(10 * time.Second); time.Since(start) < time.Second; {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			event := struct {
				From    common.Address
				Counter *big.Int
			}{}
			err := krnlContractAbi.UnpackIntoInterface(&event, "GetCounter", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Receive %v with value %v", event.From.String(), event.Counter)

			if event.Counter.Cmp(big.NewInt(10)) <= 0 {
				return errors.New("KYT_AA failed for address " + from.Hex())
			} else {
				return nil
			}
		}
	}

	return nil
}

/*
Helper method to call the AA bundler. If the response is not 200
it rejects the tx with an error.
*/
func callAABundler(payload []byte, endpoint string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Reject tx request
	if resp.StatusCode != 200 {
		log.Println("Error sending request:", err)
		log.Println(resp.Status)

		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		return nil, errors.New("Transaction rejected: request to bundler failed")
	}

	body, err := io.ReadAll(resp.Body)
	log.Println("AA Response", string(body))
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}

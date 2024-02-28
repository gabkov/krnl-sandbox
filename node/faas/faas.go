package faas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	separator := "|"
	resBody := strings.Split(txType, separator)[1]

	log.Printf("Sender %s send %s\n", sender, resBody)

	// receive the OpHash from the bundler
	txSendUserOpResponseByte, _ := callAABundler([]byte(resBody))
	txSendUserOpResponse := TxSendUserOpResponse{}
	if err := json.Unmarshal(txSendUserOpResponseByte, &txSendUserOpResponse); err != nil {
		return errors.New("Request to AA Bundler failed")
	}

	// use OpHash to get TxHash from the bundler
	txSendGetUserOpRequest := TxSendGetUserOpRequest{
		JsonRpc: "2.0",
		Id:      1,
		Method:  "eth_getUserOperationByHash",
		Params:  []string{txSendUserOpResponse.Result},
	}
	txSendGetUserOpRequestBytes, _ := json.Marshal(txSendGetUserOpRequest)
	txSendGetUserOpResponseBytes, err := callAABundler(txSendGetUserOpRequestBytes)
	txSendGetUserOpResponse := TxSendGetUserOpResponse{}
	if err := json.Unmarshal(txSendGetUserOpResponseBytes, &txSendGetUserOpResponse); err != nil {
		return errors.New("Request to AA Bundler failed")
	}

	// listen to the tx and get the response
	txHash, _ := txSendGetUserOpResponse.Params[4].(string)
	isPending := true

	for isPending {
		tx, isPending, _ := client.GetClient().TransactionByHash(context.Background(), common.BytesToHash([]byte(txHash)))
		log.Println(tx, isPending, err)

		// halt the tx base on the result
	}

	return nil
}

/*
Helper method to call the AA bundler. If the response is not 200
it rejects the tx with an error.
*/
func callAABundler(payload []byte) ([]byte, error) {
	aaBundlerEndpoint := os.Getenv("ACCOUNT_ABSTRACTION_BUNDLER_ENDPOINT")
	req, err := http.NewRequest("POST", aaBundlerEndpoint, bytes.NewBuffer(payload))
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
		return nil, errors.New("Transaction rejected: request to bundler failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}

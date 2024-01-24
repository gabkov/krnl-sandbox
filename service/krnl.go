package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gabkov/krnl-node/client"
	"io"
	"log"
	"net/http"
	"strings"
)

type RegisterDapp struct {
	DappName string `json:"dappName" binding:"required"`
}

type RegisteredDapp struct {
	AccessToken             string `json:"accessToken" binding:"required"`
	TokenAuthorityPublicKey string `json:"tokenAuthorityPublicKey" binding:"required"`
}

type TxRequest struct {
	DappName    string `json:"dappName" binding:"required"`
	AccessToken string `json:"accessToken" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

type SignatureToken struct {
	SignatureToken string `json:"signatureToken" binding:"required"`
}

type RawTransaction struct {
	RawTx string `json:"rawTx" binding:"required"`
}

type TransactionHash struct {
	TxHash string `json:"txHash" binding:"required"`
}

type Krnl struct{}

const TOKEN_AUTHORITY = "http://localhost:8080" // TODO: env

// note: probably not going to be part of the node
func (t *Krnl) RegisterNewDapp(registerDapp *RegisterDapp) RegisteredDapp {
	log.Println("RegisterNewDapp called")
	registerDappPayload, err := json.Marshal(registerDapp)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	body := callTokenAuthority("/register-dapp", registerDappPayload)

	var registeredDapp RegisteredDapp
	err = json.Unmarshal(body, &registeredDapp)
	if err != nil {
		fmt.Println("error unmarshalling response JSON:", err)
	}

	return registeredDapp
}

func (t *Krnl) TxRequest(txRequest *TxRequest) (SignatureToken, error) {
	log.Println("TxRequest called")
	txRequestPayload, err := json.Marshal(txRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	body := callTokenAuthority("/tx-request", txRequestPayload)
	if body == nil {
		return SignatureToken{}, errors.New("Transaction rejected: invalid access token") // reject
	}

	var signatureToken SignatureToken
	err = json.Unmarshal(body, &signatureToken)
	if err != nil {
		fmt.Println("error unmarshalling response JSON:", err)
	}

	return signatureToken, nil
}

func (t *Krnl) SendTx(rawTx *RawTransaction) (TransactionHash, error) {
	log.Println("SendTx called")

	client := client.GetClient()

	rawTxBytes, err := hex.DecodeString(rawTx.RawTx[2:])

	tx := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate stopping tx here
	// grabbing the requested FaaS services from the end of the input-data
	separator := "000000000000000000000000000000000000000000000000000000000000003a" // :
	res := strings.Split(hexutil.Encode(tx.Data()), separator)

	// if len is more than 1 some message is concatenated to the end of the input-data
	if len(res) > 1 {
		for i := 0; i < len(res) - 1; i++ {
			faas, err := hex.DecodeString(res[i+1])
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Requested FaaS:", string(faas))
			// do the Faas here ...
		}
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Println(err)
		return TransactionHash{}, err
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())

	return TransactionHash{TxHash: tx.Hash().Hex()}, nil
}

func callTokenAuthority(path string, payload []byte) []byte {
	req, err := http.NewRequest("POST", TOKEN_AUTHORITY+path, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	return body
}

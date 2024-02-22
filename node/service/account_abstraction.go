package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type AccountAbstraction struct{}

type SendUserOperationRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type SendUserOperationData struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string `json:"paymasterAndData"`
	Signature            string `json:"signature"`
}

var EntryPoint = os.Getenv("ACCOUNT_ABSTRACTION_ENTRY_POINT")

func (t *AccountAbstraction) SendUserOperation(tx *SendUserOperationRequest) {
	url := os.Getenv("ACCOUNT_ABSTRACTION_BUNDLER_ENDPOINT")

	// payload := SendUserOperationRequest{
	// 	Jsonrpc: "2.0",
	// 	ID:      1,
	// 	Method:  "eth_sendUserOperation",
	// }

	// opData := SendUserOperationData{}

	// payload.Params = append(payload.Params, opData)
	// payload.Params = append(payload.Params, EntryPoint)

	// jsonPayload, _ := json.Marshal(payload)

	jsonPayload, _ := json.Marshal(tx)
	log.Println("SendUserOperation")
	log.Println(string(jsonPayload))

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

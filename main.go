package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/rpc"
	gjson "github.com/gorilla/rpc/json"
	"io"
	"log"
	"net/http"
)

type RegitserDapp struct {
	DappName string `json:"dappName" binding:"required"`
}

type RegisteredDapp struct {
	AccessToken             string `json:"accessToken" binding:"required"`
	TokenAuthorityPublicKey string `json:"tokenAuthorityPublicKey" binding:"required"`
}

type TxRequest struct {
	DappName  string `json:"dappName" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type SignatureToken struct {
	SignatureToken string `json:"signatureToken" binding:"required"`
}

type KrnlTask struct{}

const TOKEN_AUTHORITY = "http://localhost:8080" // TODO: env

func main() {
	rpcServer := rpc.NewServer()

	rpcServer.RegisterCodec(gjson.NewCodec(), "application/json")

	rpcServer.RegisterService(new(KrnlTask), "")

	http.Handle("/krnl", rpcServer)

	log.Printf("Serving RPC server on port %d", 1337)

	http.ListenAndServe("localhost:1337", nil)
}

func (t *KrnlTask) RegisterNewDapp(r *http.Request, registerDapp *RegitserDapp, reply *RegisteredDapp) error {
	log.Println("RegisterNewDapp called")
	registerDappPayload, err := json.Marshal(registerDapp)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	body := callTokenAuthority("/register-dapp", registerDappPayload)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var registeredDapp RegisteredDapp
	err = json.Unmarshal(body, &registeredDapp)
	if err != nil {
		fmt.Println("error unmarshalling response JSON:", err)
		return nil
	}

	*reply = registeredDapp
	return nil
}

func (t *KrnlTask) TxRequest(r *http.Request, txRequest *TxRequest, reply *SignatureToken) error {
	log.Println("TxRequest called")
	txRequestPayload, err := json.Marshal(txRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	body := callTokenAuthority("/tx-request", txRequestPayload)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var signatureToken SignatureToken
	err = json.Unmarshal(body, &signatureToken)
	if err != nil {
		fmt.Println("error unmarshalling response JSON:", err)
		return nil
	}

	*reply = signatureToken
	return nil
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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	return body
}

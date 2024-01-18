package main

import (
	"log"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/rpc"
	gjson "github.com/gorilla/rpc/json"
)

type RegitserDapp struct {
	DappName string `json:"dappName"`
}

type RegisteredDapp struct {
	AccessToken               string `json:"accessToken"`
	TokenAuthorityPublicKey string `json:"tokenAuthorityPublicKey"`
}

type SendTx struct {
	DappName  string `json:"dappName" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type SignatureToken struct {
	SignatureToken string `json:"signatureToken" binding:"required"`
}

type KrnlTask int

const TOKEN_AUTHORITY = "http://localhost:8080" // TODO: env

func (t *KrnlTask) RegisterNewDapp(r *http.Request, registerDapp *RegitserDapp, reply *RegisteredDapp) error {

	fmt.Println("LOL")

	log.Printf("FOS")

	payloadJSON, err := json.Marshal(registerDapp)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	fmt.Println(TOKEN_AUTHORITY)

	req, err := http.NewRequest("POST", TOKEN_AUTHORITY + "/register-dapp", bytes.NewBuffer(payloadJSON))
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

	var registeredDapp RegisteredDapp
	err = json.Unmarshal(body, &registeredDapp)
	if err != nil {
		fmt.Println("error unmarshalling response JSON:", err)
		return nil
	}

	*reply = registeredDapp
	return nil
}

func main() {
	rpcServer := rpc.NewServer()

	rpcServer.RegisterCodec(gjson.NewCodec(), "application/json")

	rpcServer.RegisterService(new(KrnlTask), "")

	http.Handle("/krnl", rpcServer)

	log.Printf("Serving RPC server on port %d", 1337)

	http.ListenAndServe("localhost:1337", nil)
}
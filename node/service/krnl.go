package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type TxRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

type SignatureToken struct {
	SignatureToken string `json:"signatureToken" binding:"required"`
	Hash           string `json:"hash" binding:"required"`
}

type Krnl struct{}

/*
Forwarding the call from the client to the token-authority, 
then returning the response to the client. If there was an error, 
it rejects the tx-request.
*/
func (t *Krnl) TransactionRequest(txRequest *TxRequest) (SignatureToken, error) {
	log.Println("krnl_transactionRequest")
	txRequestPayload, err := json.Marshal(txRequest)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
	}

	body, err := callTokenAuthority("/tx-request", txRequestPayload)
	if err != nil {
		log.Println(err)
		return SignatureToken{}, errors.New(err.Error()) // reject
	}

	var signatureToken SignatureToken
	err = json.Unmarshal(body, &signatureToken)
	if err != nil {
		log.Println("error unmarshalling response JSON:", err)
	}

	return signatureToken, nil
}


/*
Helper method to call the token-auhtority api. If the response was 401
it rejects the tx with invalid access token error.
*/
func callTokenAuthority(path string, payload []byte) ([]byte, error) {
	tokenAuthority := os.Getenv("TOKEN_AUTHORITY")
	if tokenAuthority == "" {
		tokenAuthority = "http://127.0.0.1:8181" // local run
	}
	req, err := http.NewRequest("POST", tokenAuthority+path, bytes.NewBuffer(payload))
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
	if resp.StatusCode == 401 {
		return nil, errors.New("transaction rejected: invalid access token")
	}

	if resp.StatusCode == 400 {
		return nil, errors.New("transaction rejected: no FaaS request specified")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}

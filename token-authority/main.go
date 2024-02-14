package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type RegitserDapp struct {
	DappName string `json:"dappName" binding:"required"`
}

type DappSecrets struct {
	DappPK           string `json:"dappPK" binding:"required"`
	TokenAuthorityPK string `json:"tokenAuthorityPK" binding:"required"`
}

type RegisteredDapp struct {
	AccessToken             string `json:"accessToken" binding:"required"`
	TokenAuthorityPublicKey string `json:"tokenAuthorityPublicKey" binding:"required"`
}

type TxRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
	Message     string `json:"message"`
}

type SignatureToken struct {
	SignatureToken string `json:"signatureToken" binding:"required"`
	Hash           string `json:"hash" binding:"required"`
}

// mock secret
const SIGNABLE = "secret_msg"

func main() {
	router := gin.Default()

	router.POST("/register-dapp", registerDapp)
	router.POST("/tx-request", txRequest)

	router.Run(":8181")
}

/*
Generates two key-pair. One for accessing the krnl node and 
the other one as a token-authority to sing transaction requests.
Returns the accessToken and the TA public key.
*/
func registerDapp(c *gin.Context) {
	var dapp RegitserDapp
	c.BindJSON(&dapp)

	// only logging the Dapp name for now
	log.Println("Registering Dapp: ", dapp.DappName)

	dappPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dapp PK created")

	dappTokenAuthorityPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dapp TA PK created")

	dappSecrets := DappSecrets{
		DappPK:           hexutil.Encode(crypto.FromECDSA(dappPrivateKey))[2:],
		TokenAuthorityPK: hexutil.Encode(crypto.FromECDSA(dappTokenAuthorityPrivateKey))[2:]}

	file, _ := json.Marshal(dappSecrets)

	// saving the secrets into a json
	// note: with each registration new secrets being generated
	_ = os.WriteFile("secrets.json", file, 0644)

	log.Println("Dapp secrets saved")

	data := []byte(SIGNABLE)
	hash := crypto.Keccak256Hash(data)

	// creating the accessToken
	signature, err := crypto.Sign(hash.Bytes(), dappPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Access token signature created")

	taPublicKey := dappTokenAuthorityPrivateKey.Public()
	publicKeyECDSA, ok := taPublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// grabbing the TA public key which is used for signature validation 
	// in the integarting smart contracts
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	registeredDapp := RegisteredDapp{
		AccessToken:             hexutil.Encode(signature),
		TokenAuthorityPublicKey: address}

	c.JSON(200, registeredDapp)
}

/*
Validates the passed accessToken, then signs the passed message.
Returns the signature and the hash of the message.
*/
func txRequest(c *gin.Context) {
	var sendTx TxRequest
	c.BindJSON(&sendTx)

	if sendTx.Message == ""{
		log.Println("No FaaS request message provided")
		c.Status(http.StatusBadRequest)
		return
	}

	secrets, _ := os.ReadFile("secrets.json")

	dappSecrets := DappSecrets{}

	_ = json.Unmarshal([]byte(secrets), &dappSecrets)

	dappPk, err := crypto.HexToECDSA(dappSecrets.DappPK)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := dappPk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// validate accessToken
	sig, _ := hexutil.Decode(sendTx.AccessToken)
	signatureNoRecoverID := sig[:len(sig)-1] // remove recovery ID
	hash := crypto.Keccak256Hash([]byte(SIGNABLE))

	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)

	if !verified {
		log.Println("Access Token invalid")
		c.Status(http.StatusUnauthorized) // reject invalid accessToken
		return
	}

	log.Println("Access token valid")

	// create singature token
	dappTaPk, err := crypto.HexToECDSA(dappSecrets.TokenAuthorityPK)
	if err != nil {
		log.Fatal(err)
	}

	hash = crypto.Keccak256Hash([]byte(sendTx.Message))

	// sign FaaS request message
	signatureToken, err := crypto.Sign(hash.Bytes(), dappTaPk)
	if err != nil {
		log.Fatal(err)
	}

	sigToken := SignatureToken{SignatureToken: hexutil.Encode(signatureToken), Hash: hash.String()}

	c.JSON(200, sigToken)
}
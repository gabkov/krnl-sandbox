package faas

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	etrpc "github.com/ethereum/go-ethereum/rpc"
)

// simulate KYT database
// the first address from the local hardhat node config
var kytAddresses = []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}

func CallService(faas string, tx *types.Transaction) error {
	log.Println("CallService")
	log.Println(strings.TrimSpace(faas))
	switch strings.TrimSpace(faas) {
	case "KYC":
		return kyc(tx)
	case "KYT":
		return kyt(tx)
	case "KYT_AA":
		return kytAA(tx)
	default:
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

func kytAA(tx *types.Transaction) error {
	log.Println("kytAA")
	//from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	//if err != nil {
	//	log.Fatal("Could not get sender")
	//}
	//sender := from.Hex()

	//fake sender wallet
	sender := "0x8CF496044F3b5cdfdfD416a75DB9bEE798A431f7"
	from := common.HexToAddress(sender)
	// get AA request from tx data
	// KYT_AA|hex(request_body)
	salt := common.HexToHash("0x")
	accountInitCode := "0x"
	// caculate account address
	accountAddress := getCreate2Address(from, salt, accountInitCode)
	//get nonce
	entryPointAddress := common.HexToAddress(os.Getenv("STACKUP_ENTRYPOINT_ADDRESS"))
	nonce := getNonce(entryPointAddress, accountAddress, big.NewInt(3))

	// gen initCode
	factoryAddress := common.HexToAddress(os.Getenv("STACKUP_ACCOUNT_FACTORY_ADDRESS"))
	initCode := createInitCode(factoryAddress, accountAddress, big.NewInt(1))

	// call gitcoin passport

	//set userOp data
	data := map[string]interface{}{
		"sender":               accountAddress,
		"nonce":                nonce,
		"initCode":             initCode,
		"callData":             []byte(""),
		"callGasLimit":         big.NewInt(200000),
		"verificationGasLimit": big.NewInt(100000),
		"preVerificationGas":   big.NewInt(300000),
		"maxFeePerGas":         big.NewInt(1000000000),
		"maxPriorityFeePerGas": big.NewInt(500000000),
		"paymasterAndData":     []byte(""),
		"signature":            []byte(""),
	}

	// Create userOp
	op, err := userop.New(data)
	if err != nil {
		log.Fatal("Failed to create UserOperation:", err)
	}

	//get signature
	chainIDStr := os.Getenv("CHAIN_ID")
	chainID, ok := new(big.Int).SetString(chainIDStr, 10)
	if !ok {
		log.Fatal("Invalid chain ID:", chainIDStr)
	}
	signatureHash := getUserOpHash(op, entryPointAddress, chainID)
	op.Signature = signatureHash.Bytes()

	jsonData, err := op.MarshalJSON()
	if err != nil {
		log.Fatal("Failed to marshal UserOperation to JSON:", err)
	}
	fmt.Println("UserOperation as JSON:", string(jsonData))

	bunderEndpoint := ""
	callAABundler([]byte(jsonData), bunderEndpoint)

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

func getCreate2Address(fromAddress common.Address, salt common.Hash, initCode string) common.Address {
	// Convert the initCode string to bytes
	initCodeBytes, err := hex.DecodeString(initCode[2:]) // Skip the '0x' prefix
	if err != nil {
		log.Fatal("Failed to decode initCode:", err)
		return common.BytesToAddress([]byte(""))
	}

	// Calculate the hash of the initCode
	initCodeHash := sha256.Sum256(initCodeBytes)

	// Combine the bytes of the contract deployment bytecode, the salt, and the creating address
	data := make([]byte, 0, len(fromAddress.Bytes())+len(salt.Bytes())+len(initCodeHash))
	data = append(data, fromAddress.Bytes()...)
	data = append(data, salt.Bytes()...)
	data = append(data, initCodeHash[:]...)

	// Hash the result
	hashed := sha256.Sum256(data)

	// Return the address
	return common.BytesToAddress(hashed[12:])
}

func createInitCode(factoryAddress common.Address, accountAddress common.Address, key *big.Int) []byte {
	client := dialToChain()

	// ABI for factory contract
	const factoryABIJSON = `[{"inputs":[{"internalType":"contract IEntryPoint","name":"_entryPoint","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"accountImplementation","outputs":[{"internalType":"contract SimpleAccount","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[{"internalType":"contract SimpleAccount","name":"ret","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"getAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

	contractBinding := ContractBinding{
		Address: factoryAddress,
		ABI:     factoryABIJSON,
		Client:  client,
	}
	// Call GetInitCode
	initCode, err := contractBinding.GetInitCode(context.Background(), accountAddress, key)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	return initCode
}

func getNonce(contractAddress, senderAddress common.Address, key *big.Int) *big.Int {
	client := dialToChain()
	// entryPoint ABI
	abiJSON := `[{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"uint192","name":"key","type":"uint192"}],"name":"getNonce","outputs":[{"internalType":"uint256","name":"nonce","type":"uint256"}],"stateMutability":"view","type":"function"}]`

	contractBinding := ContractBinding{
		Address: contractAddress,
		ABI:     abiJSON,
		Client:  client,
	}

	// call getNonce from smart contract
	nonce, err := contractBinding.GetNonce(context.Background(), senderAddress, key)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	return nonce
}

func dialToChain() *etrpc.Client {
	client, err := etrpc.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}
	defer client.Close()
	return client
}

// ContractBinding represents a binding object for calling functions from smart contracts
type ContractBinding struct {
	Address common.Address //smart contract address
	ABI     string         // ABI
	Client  *etrpc.Client  // Client JSON-RPC
}

// GetNonce call getNonce from smart contract and return nonce
func (c *ContractBinding) GetNonce(ctx context.Context, sender common.Address, key *big.Int) (*big.Int, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return nil, err
	}

	// Encode input of getNonce
	input, err := contractABI.Pack("getNonce", sender, key)
	if err != nil {
		return nil, err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)

	nonce := new(big.Int)
	nonce, ok := nonce.SetString(result[2:], 16)
	if !ok {
		return nil, err
	}

	return nonce, nil
}

func (c *ContractBinding) GetInitCode(ctx context.Context, sender common.Address, key *big.Int) ([]byte, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return nil, err
	}

	// Encode input of createAccount
	input, err := contractABI.Pack("createAccount", sender, key)
	if err != nil {
		return nil, err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)

	initCode, err := hex.DecodeString(result[2:])
	if err != nil {
		return nil, err
	}

	return initCode, nil
}

func (c *ContractBinding) getDataFromContract(ctx context.Context, input []byte) (string, error) {
	var result string
	err := c.Client.CallContext(ctx, &result, "eth_call", map[string]interface{}{
		"to":   c.Address.Hex(),
		"data": hexutil.Encode(input),
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func getUserOpHash(op *userop.UserOperation, entryPoint common.Address, chainID *big.Int) common.Hash {
	return crypto.Keccak256Hash(
		crypto.Keccak256(op.PackForSignature()),
		common.LeftPadBytes(entryPoint.Bytes(), 32),
		common.LeftPadBytes(chainID.Bytes(), 32),
	)
}

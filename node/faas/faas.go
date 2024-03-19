package faas

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	etrpc "github.com/ethereum/go-ethereum/rpc"
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

type AABundlerParams struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// simulate KYT database
// the first address from the local hardhat node config
var kytAddresses = []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}

func CallService(faas string, tx *types.Transaction) error {
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
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}
	sender := from.Hex()
	fmt.Println("sender: ", sender)
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

	//caculate maxFeePerGas
	blockBaseFee := getMaxFeePerGas()
	newMaxFeePerGas := new(big.Int).Add(blockBaseFee, big.NewInt(1000000000))
	if newMaxFeePerGas.Cmp(blockBaseFee) < 0 {
		newMaxFeePerGas.Set(blockBaseFee)
		newMaxFeePerGas.Add(newMaxFeePerGas, big.NewInt(1))
	}

	// call gitcoin passport
	callData := []byte("0x")
	//gitCoinAddress := common.HexToAddress(os.Getenv("GITCOIN_DECODE_ADDRESS"))
	//callData := callData(gitCoinAddress, from)

	//set userOp data
	data := map[string]interface{}{
		"sender":               accountAddress,
		"nonce":                nonce,
		"initCode":             initCode,
		"callData":             callData,
		"callGasLimit":         big.NewInt(200000),
		"verificationGasLimit": big.NewInt(100000),
		"preVerificationGas":   big.NewInt(300000),
		"maxFeePerGas":         newMaxFeePerGas,
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
	signatureHash := getUserOpHash(op, entryPointAddress, tx.ChainId())
	op.Signature = signatureHash.Bytes()

	aaBundlerParams := AABundlerParams{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "eth_sendUserOperation",
		Params: []interface{}{
			op,
			os.Getenv("STACKUP_ENTRYPOINT_ADDRESS"),
		},
	}
	//convert to json
	var jsonData []byte
	jsonData, err = json.MarshalIndent(aaBundlerParams, "", "\t")
	if err != nil {
		fmt.Println("Error Marshal AABundler JSON:", err)
		return nil
	}

	fmt.Println("jsonData: ", string(jsonData))
	// receive the OpHash from the bundler
	txSendUserOpResponseBytes, err := callAABundler(jsonData, os.Getenv("ACCOUNT_ABSTRACTION_BUNDLER_ENDPOINT"))
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
		return err
	}
	krnlContractAbi, err := abi.JSON(strings.NewReader(string(krnldapp.KrnldappABI)))
	if err != nil {
		return err
	}

	_ = krnlContract

	query := ethereum.FilterQuery{
		Addresses: []common.Address{krnlAddr},
	}

	logs := make(chan types.Log)

	sub, err := sepoliaClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	// checking log for 30s
	for start := time.Now().Add(10 * time.Second); time.Since(start) < time.Second; {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			event := struct {
				From    common.Address
				Counter *big.Int
			}{}
			err := krnlContractAbi.UnpackIntoInterface(&event, "GetCounter", vLog.Data)
			if err != nil {
				return err
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

func createInitCode(factoryAddress, accountAddress common.Address, key *big.Int) []byte {
	client := dialToChain(os.Getenv("SEPOLIA_RPC_ENDPOINT"))

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

func callData(gitCoinAddress, accountAddress common.Address) []byte {
	client := dialToChain(os.Getenv("SEPOLIA_OP_GITCOIN_ENDPOINT"))
	const contractABI = `{"inputs": [{"internalType": "address","name": "_user","type": "address"}],"name": "getScore","outputs": [{"internalType": "uint256","name": "","type": "uint256"}],"stateMutability": "view","type": "function"}`
	contractBinding := ContractBinding{
		Address: gitCoinAddress,
		ABI:     contractABI,
		Client:  client,
	}
	// Call GetInitCode
	initCode, err := contractBinding.GetCallData(context.Background(), accountAddress)
	if err != nil {
		log.Fatal("Failed to get CallData:", err)
	}

	return initCode
}

func getNonce(contractAddress, senderAddress common.Address, key *big.Int) *big.Int {
	client := dialToChain(os.Getenv("SEPOLIA_RPC_ENDPOINT"))
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

func dialToChain(endpoint string) *etrpc.Client {
	client, err := etrpc.Dial(endpoint)
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
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
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

func (c *ContractBinding) GetCallData(ctx context.Context, sender common.Address) ([]byte, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return nil, err
	}

	// Encode input of createAccount
	input, err := contractABI.Pack("getScore", sender)
	if err != nil {
		return nil, err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)
	if err != nil {
		return nil, err
	}

	initCode, err := hex.DecodeString(result[2:])
	if err != nil {
		return nil, err
	}

	return initCode, nil
}

func getMaxFeePerGas() *big.Int {
	log.Println("getMaxFeePerGas: ")
	client := dialToChain(os.Getenv("SEPOLIA_RPC_ENDPOINT"))

	// get current block
	var currentBlock string
	err := client.CallContext(context.Background(), &currentBlock, "eth_blockNumber")
	if err != nil {
		log.Println(err)
		return nil
	}

	currentBlockNumber := new(big.Int)
	_, success := currentBlockNumber.SetString(strings.TrimPrefix(currentBlock, "0x"), 16)
	if !success {
		log.Fatal("Failed to parse current block number")
		return nil
	}

	// get block
	var block map[string]interface{}
	err = client.CallContext(context.Background(), &block, "eth_getBlockByNumber", fmt.Sprintf("0x%x", currentBlockNumber), true)
	if err != nil {
		log.Println(err)
		return nil
	}

	if block == nil {
		log.Fatal("Failed to retrieve the latest block")
		return nil
	}

	// get basefee
	baseFeeHex := block["baseFeePerGas"].(string)
	baseFeeBigInt, success := new(big.Int).SetString(strings.TrimPrefix(baseFeeHex, "0x"), 16)
	if !success {
		log.Fatalf("Failed to parse baseFeePerGas: %s", baseFeeHex)
	}

	return baseFeeBigInt
}

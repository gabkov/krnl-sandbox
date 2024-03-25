package faas

import (
	"bytes"
	"context"
	"crypto/ecdsa"
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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	etrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"

	"github.com/gabkov/krnl-node/client"
)

// simulate KYT database
// the first address from the local hardhat node config
var kytAddresses = []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}

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

func CallService(inputFaas string, tx *types.Transaction) error {
	faas := strings.Split(strings.TrimSpace(inputFaas), "|")
	seed := ""
	if len(faas) > 1 {
		seed = faas[1]
	}
	switch faas[0] {
	case "KYC":
		return kyc(tx)
	case "KYT":
		return kyt(tx)
	case "EL_KYT":
		return elKYT(tx)
	case "KYT_AA":
		return kytAA(tx, seed)
	case "KYT_AA_GC":
		return kytAAGC(tx)
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

// In order to use elKYT the EigenLayer AVS needs to be running
// Instructions for that can be found here: https://github.com/martonmoro/krnl-el-kyt-avs
func elKYT(tx *types.Transaction) error {
	// extract sender from the tx
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}
	// The address of the EigenLayer AVS Task Manager contract (kytTaskManager)
	kytTaskManagerAddress := common.HexToAddress("0xE6E340D132b5f46d1e472DebcD681B2aBc16e57E")

	callMsg := ethereum.CallMsg{
		To:   &kytTaskManagerAddress,
		From: from,
		// getKYTForAddress()
		Data: common.FromHex("0xaf5e8556"),
	}

	res, err := client.GetElClient().CallContract(context.Background(), callMsg, nil)

	if err != nil {
		log.Println("Error calling kyt contract: ", err)
		return err
	}

	if res[len(res)-1] == byte(1) {
		log.Println("EL KYT success for address: ", from.Hex())
		return nil
	}

	return errors.New("EL KYT failed for address " + from.Hex())
}

func kytAA(tx *types.Transaction, seed string) error {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return errors.New("Could not get sender")
	}
	sender := from.Hex()
	fmt.Println("sender: ", sender)

	//fake privateKey
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return err
	}

	num := new(big.Int)
	key, ok := num.SetString(seed, 16)
	if !ok {
		return errors.New("Seed wrong format!")
	}

	fmt.Println("privateKey: ", privateKey)
	fmt.Println("key: ", key)

	// caculate account address
	factoryAddress := common.HexToAddress(os.Getenv("STACKUP_ACCOUNT_FACTORY_ADDRESS"))

	accountAddress := getCreate2Address(factoryAddress, from, key)
	// gen initCode
	initCode := createInitCode(factoryAddress, from, key)

	//get nonce
	entryPointAddress := common.HexToAddress(os.Getenv("STACKUP_ENTRYPOINT_ADDRESS"))
	nonce := getNonce(entryPointAddress, accountAddress, big.NewInt(0))

	//caculate maxFeePerGas
	blockBaseFee := getMaxFeePerGas()
	newMaxFeePerGas := new(big.Int).Add(blockBaseFee, big.NewInt(1000000000))
	if newMaxFeePerGas.Cmp(blockBaseFee) < 0 {
		newMaxFeePerGas.Set(blockBaseFee)
		newMaxFeePerGas.Add(newMaxFeePerGas, big.NewInt(1))
	}

	// get callData
	callData := callData(entryPointAddress, from)

	//set userOp data
	data := map[string]interface{}{
		"sender":               accountAddress,
		"nonce":                nonce,
		"initCode":             initCode,
		"callData":             callData,
		"callGasLimit":         big.NewInt(6000000),
		"verificationGasLimit": big.NewInt(6000000),
		"preVerificationGas":   big.NewInt(6000000),
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

	chainId := big.NewInt(11155420) //11155420 optimism-sepolia chain id
	signatureHash := op.GetUserOpHash(entryPointAddress, chainId)
	userOpHash := signatureHash.Bytes()

	//rebuild sign hash
	signHash := signHash(userOpHash)
	signature, err := crypto.Sign(signHash.Bytes(), privateKey)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return err
	}
	//change recovery id offset from 0 | 1 -> 27|28
	signature[crypto.RecoveryIDOffset] += 27
	op.Signature = signature

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
		return err
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
	fmt.Println("txSendUserOpResponse: ", txSendUserOpResponse)

	return nil
}

func kytAAGC(tx *types.Transaction) error {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}
	sender := from.Hex()
	fmt.Println("sender: ", sender)

	//gitcoin GitcoinPassportDecoder ABI
	entryPointABI := `[{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getPassport","outputs":[{"components":[{"internalType":"string","name":"provider","type":"string"},{"internalType":"bytes32","name":"hash","type":"bytes32"},{"internalType":"uint64","name":"time","type":"uint64"},{"internalType":"uint64","name":"expirationTime","type":"uint64"}],"internalType":"struct Credential[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getScore","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"isHuman","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`

	//make client instance
	client := dialToChain(os.Getenv("SEPOLIA_RPC_ENDPOINT"))

	//set contract binding
	contractBinding := ContractBinding{
		Address: common.HexToAddress(os.Getenv("GITCOIN_PASSPORT_DECODER")),
		ABI:     entryPointABI,
		Client:  client,
	}
	//context
	ctx := context.Background()
	//get gitcoin passport
	passport, err := contractBinding.getPassport(ctx, from)
	if err != nil {
		//return err
	}
	//get gitcoin score
	score, err := contractBinding.getScore(ctx, from)
	if err != nil {
		return err
	}

	//get gitcoin isHuman
	isHuman, err := contractBinding.getIsHuman(ctx, from)
	if err != nil {
		return err
	}

	fmt.Println("passport: ", passport)
	fmt.Println("score: ", score)
	fmt.Println("isHuman: ", isHuman)
	if score.Cmp(big.NewInt(1)) < 0 || !isHuman {
		return errors.New("You do not meet the condition of gitcoin.")
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

func getCreate2Address(factoryAddress, accountAddress common.Address, key *big.Int) common.Address {
	client := dialToChain(os.Getenv("SEPOLIA_RPC_ENDPOINT"))

	entryPointABI := `[{"inputs":[{"internalType":"contract IEntryPoint","name":"_entryPoint","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"accountImplementation","outputs":[{"internalType":"contract SimpleAccount","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[{"internalType":"contract SimpleAccount","name":"ret","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"getAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

	contractBinding := ContractBinding{
		Address: factoryAddress,
		ABI:     entryPointABI,
		Client:  client,
	}

	// Call GetInitCode
	address, err := contractBinding.getAccountAddress(context.Background(), accountAddress, key)
	if err != nil {
		log.Fatal("Failed to get address:", err)
	}

	// Return the address
	return common.HexToAddress(address)
}

func createInitCode(factoryAddress, accountAddress common.Address, key *big.Int) []byte {
	// ABI for factory contract
	const factoryABIJSON = `[{"inputs":[{"internalType":"contract IEntryPoint","name":"_entryPoint","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"accountImplementation","outputs":[{"internalType":"contract SimpleAccount","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[{"internalType":"contract SimpleAccount","name":"ret","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"getAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

	contractBinding := ContractBinding{
		Address: factoryAddress,
		ABI:     factoryABIJSON,
	}
	// Call GetInitCode
	initCode, err := contractBinding.getInitCode(context.Background(), accountAddress, key)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	return initCode
}

func callData(entryPointAddress, accountAddress common.Address) []byte {
	accountABI := `[{"type":"function","name":"execute","inputs":[{"type":"address","name":"to"},{"type":"uint256","name":"value"},{"type":"bytes","name":"data"}]}]`
	contractABI := `[{"type":"function","name":"transfer","inputs":[{"type":"address","name":"to"},{"type":"uint256","name":"amount"}],"outputs":[{"type":"bool"}]}]`

	// Parse ABI
	account, err := abi.JSON(strings.NewReader(accountABI))
	if err != nil {
		fmt.Println("Error parsing account ABI:", err)
		return nil
	}
	contract, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		fmt.Println("Error parsing contract ABI:", err)
		return nil
	}

	amount := big.NewInt(0)
	inputContract, err := contract.Pack("transfer", accountAddress, amount)
	if err != nil {
		fmt.Println("err inputContract: ", err)
		return nil
	}

	inputAccount, err := account.Pack("execute", entryPointAddress, amount, inputContract)
	if err != nil {
		fmt.Println("err inputAccount: ", err)
		return nil
	}

	return inputAccount
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
	nonce, err := contractBinding.getNonce(context.Background(), senderAddress, key)
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
func (c *ContractBinding) getNonce(ctx context.Context, sender common.Address, key *big.Int) (*big.Int, error) {
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

// GetAccountAddress call getAddress from smart contract and return address
func (c *ContractBinding) getAccountAddress(ctx context.Context, owner common.Address, key *big.Int) (string, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return "", err
	}

	// Encode input of getNonce
	input, err := contractABI.Pack("getAddress", owner, key)
	if err != nil {
		return "", err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *ContractBinding) getInitCode(ctx context.Context, sender common.Address, key *big.Int) ([]byte, error) {
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

	//concat factoryAdrress with callData
	initCode := c.Address.Hex() + hex.EncodeToString(input)
	initCode = strings.TrimPrefix(initCode, "0x")
	initCodeBytes, err := hex.DecodeString(initCode)
	if err != nil {
		fmt.Println("Error decoding initCode:", err)
		return nil, err
	}

	return initCodeBytes, nil
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

func (c *ContractBinding) getCallData(ctx context.Context, sender common.Address) ([]byte, error) {
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

func getPassport(client *etrpc.Client, ctx context.Context, factoryAddress, accountAddress common.Address) common.Address {
	entryPointABI := `[{"inputs":[{"internalType":"contract IEntryPoint","name":"_entryPoint","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"accountImplementation","outputs":[{"internalType":"contract SimpleAccount","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[{"internalType":"contract SimpleAccount","name":"ret","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"getAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

	contractBinding := ContractBinding{
		Address: factoryAddress,
		ABI:     entryPointABI,
		Client:  client,
	}

	// Call GetInitCode
	address, err := contractBinding.getPassport(ctx, accountAddress)
	if err != nil {
		log.Fatal("Failed to get address:", err)
	}

	// Return the address
	return common.HexToAddress(address)
}

// getPassport call getPassport from smart contract
func (c *ContractBinding) getPassport(ctx context.Context, accountAddress common.Address) (string, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return "", err
	}

	// Encode input of getNonce
	input, err := contractABI.Pack("getPassport", accountAddress)
	if err != nil {
		return "", err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)
	if err != nil {
		return "", err
	}

	return result, nil
}

// getScore call getScore from smart contract
func (c *ContractBinding) getScore(ctx context.Context, accountAddress common.Address) (*big.Int, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return nil, err
	}

	// Encode input of getNonce
	input, err := contractABI.Pack("getScore", accountAddress)
	if err != nil {
		return nil, err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)
	if err != nil {
		return nil, err
	}

	//convert hexstring to big int
	bigInt := new(big.Int)
	score, success := bigInt.SetString(result, 0)
	if !success {
		return nil, errors.New("Can't convert hex string to big.Int")
	}

	return score, nil
}

// getIsHuman call isHuman from smart contract
func (c *ContractBinding) getIsHuman(ctx context.Context, accountAddress common.Address) (bool, error) {
	// Compile ABI
	contractABI, err := abi.JSON(strings.NewReader(c.ABI))
	if err != nil {
		return false, err
	}

	// Encode input of getNonce
	input, err := contractABI.Pack("isHuman", accountAddress)
	if err != nil {
		return false, err
	}

	// call via JSON-RPC
	var result string
	result, err = c.getDataFromContract(ctx, input)
	if err != nil {
		return false, err
	}

	//convert hexstring to big int
	bigInt := new(big.Int)
	tmpIsHuman, success := bigInt.SetString(result, 0)
	if !success {
		return false, errors.New("Can't convert hex string to big.Int")
	}

	isHuman := tmpIsHuman.Cmp(big.NewInt(0)) != 0
	return isHuman, nil
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

func getKytAAParam(tx *types.Transaction) (*ecdsa.PrivateKey, *big.Int, error) {
	separator := "000000000000000000000000000000000000000000000000000000000000003a" // :
	res := strings.Split(hexutil.Encode(tx.Data()), separator)
	if len(res) < 3 {
		return nil, nil, errors.New("Not enough parameters")
	}
	fmt.Println("res0: ", res[0])
	fmt.Println("res1: ", res[1])
	fmt.Println("res2: ", res[2])
	//fake privateKey
	tmpSeed, err := hex.DecodeString(res[2])
	if err != nil {
		fmt.Println("err1: ", err)
		return nil, nil, err
	}
	res[3] = string(bytes.Trim(tmpSeed, "\x00"))
	res[2] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	fmt.Println("res: ", res)
	// Gen ECDSA from private key
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return nil, nil, err
	}

	num := new(big.Int)
	seed, ok := num.SetString(res[3], 16)
	if !ok {
		return nil, nil, errors.New("Seed wrong format!")
	}
	fmt.Println("seed: ", res[3])

	return privateKey, seed, nil
}

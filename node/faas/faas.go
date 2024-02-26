package faas

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"os"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gabkov/krnl-node/client"
)

// simulate KYT database
// the first address from the local hardhat node config
var kytAddresses = []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}

func CallService(faas string, tx *types.Transaction) error {
	switch strings.TrimSpace(faas) {
	case "KYC":
		return kyc(tx)
	case "KYT":
		return kyt(tx)
	case "PE":
		return policyEngine(tx)
	default:
		return errors.New("unknown function name: " + faas)
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

type PEA struct {
	PolicyEngineAddress string `json:"policyEngineAddress" binding:"required"`
}

// used for lending protocol metamask demo
func policyEngine(tx *types.Transaction) error {
	pea := PEA{}
	// this contract can be deployed to any network and called there
	fileBytes, _ := os.ReadFile("./_hardhat/scripts/deployments/addresses.json")
	err := json.Unmarshal(fileBytes, &pea)
	if err != nil {
		log.Println(err)
		return err
	}
	
	policyEngineAddress := common.HexToAddress(pea.PolicyEngineAddress)
	toAddress := tx.To().String()[2:]

	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}

	callMsg := ethereum.CallMsg{
		To: &policyEngineAddress,
		From: from,
		// isAllowed(address)
		Data: common.FromHex("0xbabcc539000000000000000000000000" + toAddress),
	}

	res, err := client.GetClient().CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Println("Error calling contract: ", err)
		return err
	}

	allowed := new(big.Int)
	allowed.SetString(common.Bytes2Hex(res), 16)

	if allowed.Uint64() == 0 {
		return errors.New("policy engine - unrecognised receiver")
	}
	log.Println("Tx allowed by Policy Engine")
	return nil
}

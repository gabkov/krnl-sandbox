package faas

import (
	"context"
	"errors"
	"log"
	"math/big"
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

// used for lending protocol metamask demo
func policyEngine(tx *types.Transaction) error {
	policyEngineAddress := common.HexToAddress("0xaa292e8611adf267e563f334ee42320ac96d0463")

	toAddress := tx.To().String()[2:]

	callMsg := ethereum.CallMsg{
        To:   &policyEngineAddress,
		// isAllowed(address)
        Data: common.FromHex("0xbabcc539000000000000000000000000" + toAddress),
    }

    res, err := client.GetClient().CallContract(context.Background(), callMsg, nil)
    if err != nil {
        log.Println("Error calling contract: ", err)
    }

	allowed := new(big.Int)
	allowed.SetString(common.Bytes2Hex(res), 16)
	
	if allowed.Uint64() == 0 {
		return errors.New("unrecognised receiver")
	}
	log.Println("Tx allowed by Policy Engine")
	return nil
}

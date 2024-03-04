package faas

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
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
	case "EL_KYT":
		return elKYT(tx)
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

func elKYT(tx *types.Transaction) error {
	// extract sender from the tx
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Fatal("Could not get sender")
	}

	kytTaskManagerAddress := common.HexToAddress("0xE6E340D132b5f46d1e472DebcD681B2aBc16e57E")

	callMsg := ethereum.CallMsg{
		To: &kytTaskManagerAddress,
		From: from,
		// getKYTForAddress()
		Data: common.FromHex("0xaf5e8556"),
	}

	res, err := client.GetElClient().CallContract(context.Background(), callMsg, nil)

	if err != nil {
		log.Println("Error calling kyt contract: ", err)
		return err
	}

	if res[len(res) -1] == byte(1) {
		log.Println("EL KYT success for address: ", from.Hex())
		return nil
	}



	return errors.New("KYT failed for address " + from.Hex())
}

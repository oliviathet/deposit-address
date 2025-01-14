package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	deposit_address "github.com/lombard-finance/deposit-address"
)

// In top level project folder, run `go run demo/cmd/derive/main.go`.
func main() {
	// This is the public key for the base deposit key on Cubist that all deposit addresses are derived from
	basePubKey := "0x043dcf7a68429b23a0396ca61c1ab243ccbbcc629ff04c59394458d6db5dd2bb159e0b7a71ef07247b59a0a21b1f1eaee61a40064ade423e926f38550065a43587"

	toAddress := "<ETH ADDRESS>"
	referral := "lombard"
	nonce := uint32(0)

	auxData, err := deposit_address.ComputeAuxDataV0(nonce, []byte(referral))
	if err != nil {
		fmt.Println("Error computing aux data:", err)
		return
	}

	// parse public key of base deposit key
	pkBytes, err := hex.DecodeString(strings.TrimPrefix(basePubKey, "0x"))
	if err != nil {
		fmt.Println("Decoding public key hex:", err)
		return
	}

	pk, err := secp256k1.ParsePubKey(pkBytes)
	if err != nil {
		log.Println("Error decoding SEC1 encoded pubkey:", err)
		return
	}

	// for ethereum mainnet, LBTC contract address is 0x8236a87084f8B84306f72007F36F2618A5634494
	lbtcContractAddrEthMainnet := [20]byte{0x82, 0x36, 0xa8, 0x70, 0x84, 0xf8, 0xB8, 0x43, 0x06, 0xf7, 0x20, 0x07, 0xF3, 0x6F, 0x26, 0x18, 0xA5, 0x63, 0x44, 0x94}
	lbtcContractAddr := common.BytesToAddress(lbtcContractAddrEthMainnet[:20])

	// for ethereum mainnet, chain id is 1
	chainId := [32]byte{}
	chainId[31] = 1

	// derive the address
	derivedAddr, err := deposit_address.EvmDepositSegwitAddr(pk, lbtcContractAddr, common.HexToAddress(toAddress), chainId[:], auxData[:], &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error tweaking address:", err)
		return
	}

	fmt.Println("Address:", derivedAddr)

}

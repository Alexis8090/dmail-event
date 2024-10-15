package main

import (
	"fmt"
	"log"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func main() {
	// Connect to Substrate node (replace with your node's WebSocket endpoint)
	// wss://archive.testnet.analog.one
	// api, err := gsrpc.NewSubstrateAPI("wss://rpc.polkadot.io")
	api, err := gsrpc.NewSubstrateAPI("wss://archive.testnet.analog.one")

	if err != nil {
		log.Fatal(err)
	}

	// Replace with your block hash (must be in hex format, starting with '0x')
	// 2,400,332  0x55ba66d64019d06d4140cb5eb1f589997808e1d721a4797a0a4119d6d9f5e42e
	// blockHash, err := types.NewHashFromHexString("0xf963acedd1748ae02b05c00eb5c1f97d16b2d32d89391a7d7ccc276eaff423a0")
	// blockHash, err := types.NewHashFromHexString("0x55ba66d64019d06d4140cb5eb1f589997808e1d721a4797a0a4119d6d9f5e42e")
	blockHash, err := api.RPC.Chain.GetBlockHash(2400332)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the block
	block, err := api.RPC.Chain.GetBlock(blockHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Block Number: %d \n", block.Block.Header.Number)

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the events associated with this block
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		log.Fatal(err)
	}

	raw, err := api.RPC.State.GetStorageRaw(key, blockHash)
	if err != nil {
		log.Fatal(err)
	}

	var events types.EventRecords
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)
	if err != nil {
		log.Fatal(err)
	}

	// // Iterate over the events and print them
	fmt.Printf("Events for block hash %s \n", blockHash.Hex())

	// fmt.Println("events as ", events.Balances_Deposit)
	for _, e := range events.Dmail_Message {
		fmt.Printf("Event: %+v\n", e.Message)
	}
}
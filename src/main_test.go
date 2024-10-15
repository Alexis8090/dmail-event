// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestChain_SubscribeFinalizedHeads(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	api, err := gsrpc.NewSubstrateAPI("wss://archive.mainnet.analog.one")
	assert.NoError(t, err)

	sub, err := api.RPC.Chain.SubscribeFinalizedHeads()
	assert.NoError(t, err)
	defer sub.Unsubscribe()


	for {
		select {
		case head := <-sub.Chan():
			blockHash, err := api.RPC.Chain.GetBlockHash(uint64(head.Number))
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
	}
}

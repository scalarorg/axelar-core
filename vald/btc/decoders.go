package btc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"

	"github.com/scalarorg/btc-vault/btcvault"
)

func DecodeEventContractCall(tx *rpc.BTCTransaction) (types.EventContractCall, error) {
	fmt.Println(btcvault.TagLen)
	// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
	sender := types.Address(common.BytesToAddress([]byte(tx.Data.Hash)))
	destinationChain := nexus.ChainName("ethereum-sepolia")
	contractAddress := "0xe432150cce91c13a887f7D836923d5597adD8E31"

	return types.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      types.Hash([]byte(tx.Data.Hash)),
	}, nil
}

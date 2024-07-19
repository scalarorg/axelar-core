package btc

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
)

func DecodeEventContractCall(tx *rpc.BTCTransaction) (types.EventContractCall, error) {
	// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
	sender := types.Address(common.BytesToAddress([]byte(tx.Metadata.Details[0].Address)))
	destinationChain := nexus.ChainName("Ethereum")
	contractAddress := ""

	return types.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      types.Hash(tx.RawData.Hash().CloneBytes()),
	}, nil
}

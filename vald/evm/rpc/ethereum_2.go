package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// Ethereum2Client is a JSON-RPC client of Ethereum 2.0
type Ethereum2Client struct {
	*EthereumClient
}

// NewEthereum2Client is the constructor
func NewEthereum2Client(ethereumClient *EthereumClient) (*Ethereum2Client, error) {
	client := &Ethereum2Client{EthereumClient: ethereumClient}
	if _, err := client.LatestFinalizedBlockNumber(context.Background(), 0); err != nil {
		fmt.Println("error getting finalized block number from ethereum_2 - eth-debugging", err)
		return nil, err
	}

	return client, nil
}

// LatestFinalizedBlockNumber returns the latest finalized block number
func (c *Ethereum2Client) LatestFinalizedBlockNumber(ctx context.Context, _ uint64) (*big.Int, error) {
	var head *types.Header
	err := c.rpc.CallContext(ctx, &head, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		fmt.Println("error getting finalized block number from ethereum_2 - eth-debugging", err)
		return nil, err
	}
	if head == nil || head.Number == nil {
		fmt.Println("error getting finalized block number from ethereum_2 - eth-debugging", ethereum.NotFound)
		return nil, ethereum.NotFound
	}

	fmt.Println("eth-debugging - newBlockHeight from ethereum_2: ", head.Number)

	return head.Number, nil
}

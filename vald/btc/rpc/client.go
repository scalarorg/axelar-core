package rpc

import (
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/btcsuite/btcd/rpcclient"
)

type TxResult results.Result[BTCTransaction]

type Client interface {
	GetTransaction(txID types.Hash) (BTCTransaction, error)
	GetTransactions(txIDs []types.Hash) ([]TxResult, error)
	Close()
}

func NewClient(cfg *rpcclient.ConnConfig, logger log.Logger) (Client, error) {
	client, error := NewBTCClient(cfg, logger)

	if error != nil {
		return nil, error
	}

	return client, nil
}

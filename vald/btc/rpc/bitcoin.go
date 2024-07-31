package rpc

import (
	"sync"

	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type BTCClient struct {
	client *rpcclient.Client
	logger log.Logger
	cfg    *rpcclient.ConnConfig
}

func NewBTCClient(cfg *rpcclient.ConnConfig, logger log.Logger) (*BTCClient, error) {
	client, err := rpcclient.New(cfg, nil)
	if err != nil {
		return nil, err
	}

	return &BTCClient{
		client: client,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (c *BTCClient) Close() {
	c.client.Shutdown()
}

type BTCTransaction struct {
	Data btcjson.TxRawResult
}

func (c *BTCClient) GetTransaction(txID types.Hash) (BTCTransaction, error) {
	var tx BTCTransaction
	txBytes := txID.Bytes()
	txBytesReverse := make([]byte, len(txBytes))
	for i, b := range txBytes {
		txBytesReverse[len(txBytes)-1-i] = b
	}

	txHash, err := chainhash.NewHash(txBytesReverse)

	if err != nil {
		c.logger.Errorf("failed to create BTC chainhash from txID", "txID", txID, "error", err)
		return tx, err
	}

	txMetadata, err := c.client.GetRawTransactionVerbose(txHash)

	if err != nil {
		c.logger.Errorf("failed to get BTC transaction", "txID", txID, "error", err)
	} else {
		tx.Data = *txMetadata
	}

	return tx, err
}

func (c *BTCClient) GetTransactions(txIDs []types.Hash) ([]TxResult, error) {
	txs := make([]TxResult, len(txIDs))
	var wg sync.WaitGroup

	for i := range txIDs {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			tx, err := c.GetTransaction(txIDs[index])

			var txResult TxResult
			if err != nil {
				txResult = TxResult(results.FromErr[BTCTransaction](err))
			} else {
				txResult = TxResult(results.FromOk[BTCTransaction](tx))
			}

			txs[index] = txResult

		}(i)
	}

	wg.Wait()
	return txs, nil
}

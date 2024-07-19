package rpc

import (
	"sync"

	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
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
	Metadata btcjson.GetTransactionResult
	RawData  btcutil.Tx
}

func (c *BTCClient) GetTransaction(txID types.Hash) (BTCTransaction, error) {
	var tx BTCTransaction
	txHash, err := chainhash.NewHash(txID.Bytes())
	if err != nil {
		c.logger.Errorf("failed to create BTC chainhash from txID", "txID", txID, "error", err)
		return tx, err
	}

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		wg.Add(1)
		txMetadata, err := c.client.GetTransaction(txHash)

		if err != nil {
			c.logger.Errorf("failed to get BTC transaction", "txID", txID, "error", err)
		} else {
			tx.Metadata = *txMetadata
		}
	}()

	go func() {
		txDetail, err := c.client.GetRawTransaction(txHash)

		if err != nil {
			c.logger.Errorf("failed to get BTC raw transaction", "txID", txID, "error", err)

		} else {
			tx.RawData = *txDetail
		}
	}()

	wg.Wait()

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

package btc

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/scalarorg/btc-vault/btcvault"

	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

// hard code
const (
	covenant_quorum = 3
)

var net = &chaincfg.MainNetParams

var expected_tag = []byte("01020304")

var covenant_pks = [][]byte{
	[]byte("02a60c2b5524ee762c3e56ea23a992953b88226463abcd1ca819d1d64880393791"),
	[]byte("02f8812f5b3456f55b24d1e038a3900eae89a46d2394825b55145c24bb41c5be6d"),
	[]byte("0330c74896d8724ef941d9fb129ea2a36b2e4eb87d37299ebd2f0f61be5a8d88cf"),
	[]byte("03b029725f0b73a63cbf43cb7219178af098cd28f49398be587c75c34d5ae60fa8"),
	[]byte("03fa7eb611fd1466820d97fc682a35e0f3af5db8e1a35acf653fa2e2f482f1cc52"),
}

func covenant_pks_to_pubkeys() []*btcec.PublicKey {
	pubkeys := make([]*btcec.PublicKey, len(covenant_pks))
	for i, pk := range covenant_pks {
		pubkeys[i], _ = btcec.ParsePubKey(pk)
	}
	return pubkeys
}

// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
func DecodeEventContractCall(tx *rpc.BTCTransaction) (types.EventContractCall, error) {
	txHex := tx.Data.Hex
	// Decode the hex string into bytes
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	// Parse the transaction
	neededParseTx := wire.NewMsgTx(wire.TxVersion)
	err = neededParseTx.Deserialize(hex.NewDecoder(bytes.NewReader(txBytes)))
	if err != nil {
		log.Fatalf("failed to parse transaction: %v", err)
	}
	parseData, err := btcvault.ParseV0VaultTx(neededParseTx, expected_tag, covenant_pks_to_pubkeys(), covenant_quorum, net)
	if err != nil {
		log.Fatalf("failed to parse vault tx: %v", err)
	}
	payloadData := parseData.PayloadOpReturnData
	sender := types.Address(common.BytesToAddress(payloadData.ChainIdUserAddress))
	// destinationChain := payloadData.ChainID
	destinationChain := nexus.ChainName("ethereum-sepolia")
	contractAddress := hex.EncodeToString(payloadData.ChainIdSmartContractAddress)
	// need "0x"?

	return types.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      types.Hash([]byte(tx.Data.Hash)),
	}, nil
}

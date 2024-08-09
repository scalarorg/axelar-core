package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/scalarorg/btc-vault/btcvault"

	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type globalParams struct {
	expected_tag    string
	covenant_pks    []string
	covenant_quorum uint32
}

// hard code
const (
	path                = "params/"
	global_params_path  = path + "global_params.json"
	mainnet_params_path = path + "main_net_params.json"
	testnet_params_path = path + "test_net_params.json"
)

var net = &chaincfg.MainNetParams // it can be anything cause we don't need to compare it

// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
func DecodeEventContractCall(tx *rpc.BTCTransaction) (types.EventContractCall, error) {
	// get global params to parse the vault transactions
	var globalParams globalParams
	parseJsonFile(global_params_path, &globalParams)
	expected_tag, err := hex.DecodeString(globalParams.expected_tag)
	if err != nil {
		log.Fatalf("failed to decode hex string")
		return types.EventContractCall{}, err
	}
	covenant_pks := [][]byte{}
	for i, key := range globalParams.covenant_pks {
		covenant_pks[i], err = hex.DecodeString(key)
		if err != nil {
			log.Fatalf("failed to decode hex string")
			return types.EventContractCall{}, err
		}
	}
	covenant_quorum := globalParams.covenant_quorum

	// get mainnet or testnet chain params to detect chain
	testnetParams := make(map[string]string)
	parseJsonFile(testnet_params_path, &testnetParams)

	// Begin to parse the transaction
	txHex := tx.Data.Hex
	// Decode the hex string into bytes
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		log.Fatalf("failed to decode hex string")
		return types.EventContractCall{}, err
	}

	// Parse the transaction
	neededParseTx := wire.NewMsgTx(wire.TxVersion)
	err = neededParseTx.Deserialize(hex.NewDecoder(bytes.NewReader(txBytes)))
	if err != nil {
		log.Fatalf("failed to parse transaction")
		return types.EventContractCall{}, err
	}
	parseData, err := btcvault.ParseV0VaultTx(neededParseTx, expected_tag, covenantPksToPubkeys(covenant_pks), covenant_quorum, net)
	if err != nil {
		log.Fatalf("failed to parse vault tx:")
		return types.EventContractCall{}, err
	}
	// Get Payload data
	payloadData := parseData.PayloadOpReturnData
	// Get the sender
	sender := types.Address(common.BytesToAddress(payloadData.ChainIdUserAddress))
	// Find and Get the chain name
	byteChainID := payloadData.ChainID
	numberChainID := binary.BigEndian.Uint64(byteChainID)
	nameChainID := testnetParams[strconv.FormatUint(numberChainID, 10)]
	destinationChain := nexus.ChainName(nameChainID)
	// Get the contract address
	contractAddress := hex.EncodeToString(payloadData.ChainIdSmartContractAddress)
	// need "0x"?

	return types.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      types.Hash([]byte(tx.Data.Txid)),
	}, nil
}

func parseJsonFile(path_file string, v interface{}) {
	jsonFile, err := os.Open(path_file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer jsonFile.Close()
	bytesFile, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	err = json.Unmarshal(bytesFile, v)
	if err != nil {
		log.Fatalf("failed to parse json file: %v", err)
	}
}

func covenantPksToPubkeys(covenant_pks [][]byte) []*btcec.PublicKey {
	pubkeys := make([]*btcec.PublicKey, len(covenant_pks))
	for i, pk := range covenant_pks {
		pubkeys[i], _ = btcec.ParsePubKey(pk)
	}
	return pubkeys
}

package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/scalarorg/btc-vault/btcvault"
)

func DecodeEventContractCall(tx *rpc.BTCTransaction) (types.EventContractCall, error) {
	// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
	// sender := types.Address(common.BytesToAddress([]byte(tx.Data.Hash)))
	// destinationChain := nexus.ChainName("ethereum-sepolia")
	// contractAddress := "0xe432150cce91c13a887f7D836923d5597adD8E31"

	// Decode the hex string into bytes
	txRaw, err := hex.DecodeString(tx.Data.Hex)
	if err != nil {
		log.Fatalf("Failed to decode hex string %v", err)
		return types.EventContractCall{}, err
	}
	fmt.Printf("txRaw: %v\n", txRaw)
	// Parse the transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)
	err = msgTx.Deserialize(bytes.NewReader(txRaw))
	if err != nil {
		log.Fatalf("Failed to parse transaction")
		return types.EventContractCall{}, err
	}
	if len(msgTx.TxOut) < 3 {
		return types.EventContractCall{}, fmt.Errorf("btcLocking tx must have at least 3 outputs")
	}
	payloadData, err := btcvault.NewPayloadOpReturnDataFromTxOutput(msgTx.TxOut[2])
	if err != nil {
		return types.EventContractCall{}, fmt.Errorf("cannot parse payload op return data: %w", err)
	}

	if payloadData == nil {
		return types.EventContractCall{}, fmt.Errorf("transaction does not have expected payload op return output")
	}
	sender := types.Address(common.BytesToAddress(payloadData.ChainIdUserAddress))
	// Find and Get the chain name
	// numberChainID := binary.BigEndian.Uint64(payloadData.ChainID)
	// destinationChain := nexus.ChainName(strconv.FormatUint(numberChainID, 10))
	// Todo: Scalar hardcoded chain name for now
	destinationChain := nexus.ChainName("ethereum-sepolia")
	// Get the contract address
	contractAddress := hex.EncodeToString(payloadData.ChainIdSmartContractAddress)
	// need "0x"?

	abi_minting_payload := types.Hash(common.BytesToHash(payloadData.Amount))
	abi_address_payload := types.Hash(common.BytesToHash(payloadData.ChainIdUserAddress))

	abi_payload := append(abi_address_payload[:], abi_minting_payload[:]...)
	// Get the payload hash
	payloadHash := types.Hash(common.BytesToHash(crypto.Keccak256(abi_payload)))

	return types.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      payloadHash,
	}, nil
}

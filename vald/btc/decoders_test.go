package btc

import (
	"fmt"
	"testing"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/assert"
)

func TestBtcTransactionDecode(t *testing.T) {
	//txId := "2f8c90c7596883774127760afd2c5734f637316afe04576a69781152bc6ebf2a"
	tx := rpc.BTCTransaction{
		Data: btcjson.TxRawResult{
			Hex: "020000000001019369c111202459015d0a4e5a4b7cca594d658dce0561cfadcb338ae9f53806e60300000000fdffffff0410270000000000002251207f99d0801267696850236ed8a63bd386e151e4f5704c64ab070aa5e87299be910000000000000000476a4501020304002b122fd36a9db2698c39fb47df3e3fa615e70e368acb874ec5494e4236722b2d61e1436122e3973468bd8776b8ca0645e37a5760c4a2be7796acb94cf312ce0d00000000000000003a6a380000000000aa36a7130c4810d57140e1e62967cbf742caeae91b6ece768e8de8cf0c7747d41f75f83c914a19c5921cf30000000000002710757f00000000000016001408b7b00b0f720cf5cc3e7e38aaae1a572b962b24024830450221008cf6c01306f9a1e6c4209fbbe40bc5c9992abb1cd86ae383db468700ed9dc8d4022000d6abe69a80abbdfe6dcf17aa3b2f411c7bbd9002e1e025dcfe3f9747e532ac0121032b122fd36a9db2698c39fb47df3e3fa615e70e368acb874ec5494e4236722b2d00000000",
		},
	}
	event, err := DecodeEventContractCall(&tx)
	if err != nil {
		fmt.Printf("Error decoding event: %v", err)
		// t.Errorf("Error decoding event: %v", err)
	}
	t.Logf("Event: %v", event)
	assert.Equal(t, "11155111", event.DestinationChain.String(), "Chain ID should be 11155111")
	assert.Equal(t, "768e8de8cf0c7747d41f75f83c914a19c5921cf3", event.ContractAddress, "Chain ID should be 11155111")
}

package evm

import (
	"context"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	btc "github.com/axelarnetwork/axelar-core/vald/btc"
	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	voteTypes "github.com/axelarnetwork/axelar-core/x/vote/types"
)

// ProcessGatewayTxConfirmation votes on the correctness of an EVM chain gateway's transactions
func (mgr Mgr) ProcessGatewayTxConfirmationBTC(event *types.ConfirmGatewayTxStarted) error {
	fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 1")
	if !mgr.isParticipantOf(event.Participants) {
		mgr.logger("pollID", event.PollID).Debug("ignoring gateway tx confirmation poll: not a participant")
		return nil
	}

	fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 2")

	var vote *voteTypes.VoteRequest

	txReceipt, err := mgr.btcMgr.GetTxIfFinalized(event.TxID, event.ConfirmationHeight)

	fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 3")

	if err != nil {
		return err
	}

	if txReceipt.Err() != nil {
		fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 4")
		vote = voteTypes.NewVoteRequest(mgr.proxy, event.PollID, types.NewVoteEvents(event.Chain))

		mgr.logger().Infof("broadcasting empty vote for poll %s: %s", event.PollID.String(), txReceipt.Err().Error())
	} else {
		fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 5")
		events := mgr.processGatewayTxBTC(event.Chain, event.GatewayAddress, txReceipt.Ok(), event.TxID)
		vote = voteTypes.NewVoteRequest(mgr.proxy, event.PollID, types.NewVoteEvents(event.Chain, events...))

		mgr.logger().Infof("broadcasting vote %v for poll %s", events, event.PollID.String())
	}

	fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 6")

	_, err = mgr.broadcaster.Broadcast(context.TODO(), vote)

	fmt.Println("debugging - ProcessGatewayTxConfirmationBTC 7")

	return err
}

// extract receipt processing from ProcessGatewayTxConfirmation, so that it can be used in ProcessGatewayTxsConfirmation
func (mgr Mgr) processGatewayTxBTC(chain nexus.ChainName, _ types.Address, tx rpc.BTCTransaction, txID types.Hash) []types.Event {
	var events []types.Event

	// Temporary workaround with Event_ContractCall for BTC
	btcEvent, err := btc.DecodeEventContractCall(&tx)

	if err != nil {
		mgr.logger().Debug(sdkerrors.Wrap(err, "decode event ContractCall failed").Error())
	}

	if err := btcEvent.ValidateBasic(); err != nil {
		mgr.logger().Debug(sdkerrors.Wrap(err, "invalid event ContractCall").Error())
	}

	events = append(events, types.Event{
		Chain: chain,
		TxID:  types.Hash(txID),
		Index: 0,
		Event: &types.Event_ContractCall{
			ContractCall: &btcEvent,
		},
	})

	return events
}

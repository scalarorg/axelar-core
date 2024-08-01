package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"

	"github.com/axelarnetwork/axelar-core/x/multisig/exported"
)

var _ sdk.Msg = &StartKeygenRequest{}

// 20240801 Taivv: For amino support.
var _ legacytx.LegacyMsg = &StartKeygenRequest{}

// NewStartKeygenRequest constructor for StartKeygenRequest
func NewStartKeygenRequest(sender sdk.AccAddress, keyID exported.KeyID) *StartKeygenRequest {
	return &StartKeygenRequest{
		Sender: sender,
		KeyID:  keyID,
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (m StartKeygenRequest) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, sdkerrors.Wrap(err, "sender").Error())
	}

	if err := m.KeyID.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// GetSigners implements the sdk.Msg interface
func (m StartKeygenRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

/*
* 20240801: Scalar
* Implement the LegacyMsg interface
 */

// GetSignBytes returns the bytes all expected signers must sign over for a
// StartKeygenRequest.
func (msg StartKeygenRequest) GetSignBytes() []byte {
	var amino = codec.NewLegacyAmino()
	return sdk.MustSortJSON(amino.MustMarshalJSON(&msg))
}

func (msg StartKeygenRequest) Route() string {
	return "keygenStart"
}

func (msg StartKeygenRequest) Type() string {
	return "StartKeygenRequest"
}

/* End of LegacyMsg interface */

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

const (
	flagMultisig        = "multisig"
	flagOverwrite       = "overwrite"
	flagSigOnly         = "signature-only"
	flagAmino           = "amino"
	flagNoAutoIncrement = "no-auto-increment"
)

// Cmd creates a main CLI command
func ScalarCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scalar",
		Short: "Tool for helping with debugging your application",
		RunE:  client.ValidateCmd,
	}
	cmd.AddCommand(Debug())
	// cmd.AddCommand(PubkeyCmd())
	// cmd.AddCommand(AddrCmd())
	// cmd.AddCommand(RawBytesCmd())

	return cmd
}

// getPubKeyFromString decodes SDK PubKey using JSON marshaler.
func getPubKeyFromString(ctx client.Context, pkstr string) (cryptotypes.PubKey, error) {
	var pk cryptotypes.PubKey
	err := ctx.Codec.UnmarshalInterfaceJSON([]byte(pkstr), &pk)
	return pk, err
}

func Debug() *cobra.Command {
	return &cobra.Command{
		Use:   "debug [message]",
		Short: "Debug unsigned legacy message",
		Long: fmt.Sprintf(`Debug unsigned legacy message.

Example:
$ %s debug /root/.axelar/unsigned_msg.json`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("Scalar debug# LegacyMsg when using amino JSON")
			clientCtx, _ := client.GetClientTxContext(cmd)
			stdTx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			// f := cmd.Flags()
			// txFactory := tx.NewFactoryCLI(clientCtx, f.Flags())
			txCfg := clientCtx.TxConfig
			txBuilder, err := txCfg.WrapTxBuilder(stdTx)
			if err != nil {
				return err
			}
			msgs := txBuilder.GetTx().GetMsgs()
			msgsBytes := make([]json.RawMessage, 0, len(msgs))
			for _, msg := range msgs {
				cmd.Println("Scalar debug# Message:", msg, "#")
				legacyMsg, ok := msg.(legacytx.LegacyMsg)
				if !ok {
					cmd.Println("Scalar debug# Error:", ok, "#")
					panic(fmt.Errorf("expected %T when using amino JSON", (*legacytx.LegacyMsg)(nil)))
				}

				msgsBytes = append(msgsBytes, json.RawMessage(legacyMsg.GetSignBytes()))
			}
			return nil
		},
	}
}

func PubkeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pubkey [pubkey]",
		Short: "Decode a pubkey from proto JSON",
		Long: fmt.Sprintf(`Decode a pubkey from proto JSON and display it's address.

Example:
$ %s debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			pk, err := getPubKeyFromString(clientCtx, args[0])
			if err != nil {
				return err
			}
			cmd.Println("Address:", pk.Address())
			cmd.Println("PubKey Hex:", hex.EncodeToString(pk.Bytes()))
			return nil
		},
	}
}

func AddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Long: fmt.Sprintf(`Convert an address between hex encoding and bech32.

Example:
$ %s debug addr cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrString := args[0]
			var addr []byte

			// try hex, then bech32
			var err error
			addr, err = hex.DecodeString(addrString)
			if err != nil {
				var err2 error
				addr, err2 = sdk.AccAddressFromBech32(addrString)
				if err2 != nil {
					var err3 error
					addr, err3 = sdk.ValAddressFromBech32(addrString)

					if err3 != nil {
						return fmt.Errorf("expected hex or bech32. Got errors: hex: %v, bech32 acc: %v, bech32 val: %v", err, err2, err3)
					}
				}
			}

			cmd.Println("Address:", addr)
			cmd.Printf("Address (hex): %X\n", addr)
			cmd.Printf("Bech32 Acc: %s\n", sdk.AccAddress(addr))
			cmd.Printf("Bech32 Val: %s\n", sdk.ValAddress(addr))
			return nil
		},
	}
}

func RawBytesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "raw-bytes [raw-bytes]",
		Short: "Convert raw bytes output (eg. [10 21 13 255]) to hex",
		Long: fmt.Sprintf(`Convert raw-bytes to hex.

Example:
$ %s debug raw-bytes [72 101 108 108 111 44 32 112 108 97 121 103 114 111 117 110 100]
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			stringBytes := args[0]
			stringBytes = strings.Trim(stringBytes, "[")
			stringBytes = strings.Trim(stringBytes, "]")
			spl := strings.Split(stringBytes, " ")

			byteArray := []byte{}
			for _, s := range spl {
				b, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return err
				}
				byteArray = append(byteArray, byte(b))
			}
			fmt.Printf("%X\n", byteArray)
			return nil
		},
	}
}

package cli

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	eciesgo "github.com/ecies/go/v2"
	"github.com/jackal-dao/canine/x/filetree/types"
	filetypes "github.com/jackal-dao/canine/x/filetree/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPostFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post-file [path] [account] [contents] [keys] [viewers] [editors]",
		Short: "post a new file to an account's file explorer",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argHashpath := args[0]
			argAccount := args[1]
			argContents := args[2]
			argKeys := args[3]
			argViewers := args[4]
			argEditors := args[5]

			viewerAddresses := strings.Split(argViewers, ",")
			editorAddresses := strings.Split(argEditors, ",")

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//Cut out the / at the end to save heartache
			trimPath := strings.TrimSuffix(argHashpath, "/")
			chunks := strings.Split(trimPath, "/")

			//The problem: In order to check editor or viewer access, access.go functions need to consume a file, and they hex( hash (file.address, user)) to check if the user
			//In question has viewing access. Unfortunately, unless we build the full merklePath on the CLI side and pass it into our editors/viewers address map, our access.go
			//functions don't work. I guess we could build the viewers and editors inside of the keeper, but that would be so expensive.
			//However, inside of the keeper, we do a check to ensure that the caller of postFile msg is in fact the owner of the parent file...so if
			//They are not, then the JSON viewers and editors that we just passed in, won't even be used because the child file will never be created!
			fullMerklePath := types.MerklePath(trimPath)

			parentString := strings.Join(chunks[0:len(chunks)-1], "/")
			fmt.Println("parentString is", parentString)

			childString := string(chunks[len(chunks)-1])
			fmt.Println("ChildString is:", childString)

			parentHash := types.MerklePath(parentString)
			fmt.Println("parent hash is", parentHash)

			h := sha256.New()
			h.Write([]byte(childString))
			childHash := fmt.Sprintf("%x", h.Sum(nil))
			fmt.Println("child hash is", childHash)

			viewers := make(map[string]string)
			editors := make(map[string]string)

			viewerAddresses = append(viewerAddresses, clientCtx.GetFromAddress().String())
			editorAddresses = append(editorAddresses, clientCtx.GetFromAddress().String())

			for _, v := range viewerAddresses {
				if len(v) < 1 {
					continue
				}
				key, err := sdk.AccAddressFromBech32(v)
				if err != nil {
					return err
				}
				//As of right now, if you post a file and try to add viewers that have not initiated an account, this will error and say "public key not found, perhaps not inited yet"
				queryClient := filetypes.NewQueryClient(clientCtx)
				res, err := queryClient.Pubkey(cmd.Context(), &filetypes.QueryGetPubkeyRequest{Address: key.String()})
				if err != nil {
					return types.ErrPubKeyNotFound
				}

				pkey, err := eciesgo.NewPublicKeyFromHex(res.Pubkey.Key)
				if err != nil {
					return err
				}

				encrypted, err := clientCtx.Keyring.Encrypt(pkey.Bytes(false), []byte(argKeys))
				if err != nil {
					return err
				}

				h := sha256.New()
				h.Write([]byte(fmt.Sprintf("v%s%s", fullMerklePath, v)))
				hash := h.Sum(nil)

				addressString := fmt.Sprintf("%x", hash)

				viewers[addressString] = fmt.Sprintf("%x", encrypted)
			}

			for _, v := range editorAddresses {
				if len(v) < 1 {
					continue
				}
				fmt.Println(v)

				key, err := sdk.AccAddressFromBech32(v)
				if err != nil {
					return err
				}

				queryClient := authtypes.NewQueryClient(clientCtx)
				res, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{Address: key.String()})
				if err != nil {
					return err
				}

				var acc authtypes.BaseAccount

				err = acc.Unmarshal(res.Account.Value)
				if err != nil {
					return err
				}
				var pkey secp256k1.PubKey

				err = pkey.Unmarshal(acc.PubKey.Value)
				if err != nil {
					return err
				}

				encrypted, err := clientCtx.Keyring.Encrypt(pkey.Key, []byte(argKeys))
				if err != nil {
					return err
				}

				h := sha256.New()
				h.Write([]byte(fmt.Sprintf("e%s%s", fullMerklePath, v))) //this used to be pathString
				hash := h.Sum(nil)

				addressString := fmt.Sprintf("%x", hash)

				editors[addressString] = fmt.Sprintf("%x", encrypted)
			}

			jsonViewers, err := json.Marshal(viewers)
			if err != nil {
				return err
			}

			jsonEditors, err := json.Marshal(editors)
			if err != nil {
				return err
			}

			msg := types.NewMsgPostFile(
				clientCtx.GetFromAddress().String(), //Sender of msg
				argAccount,
				parentHash,
				childHash,
				argContents, //child piece
				string(jsonViewers),
				string(jsonEditors),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

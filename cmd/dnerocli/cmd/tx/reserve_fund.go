package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	"github.com/dnerochain/dnero/ledger/types"
	"github.com/dnerochain/dnero/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// reserveFundCmd represents the reserve fund command
// Example:
//		dnerocli tx reserve --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --fund=900 --collateral=1203 --seq=6 --duration=1002 --resource_ids=die_another_day,hello
var reserveFundCmd = &cobra.Command{
	Use:     "reserve",
	Short:   "Reserve fund for an off-chain micropayment",
	Example: `dnerocli tx reserve --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --fund=900 --collateral=1203 --seq=6 --duration=1002 --resource_ids=die_another_day,hello`,
	Run:     doReserveFundCmd,
}

func doReserveFundCmd(cmd *cobra.Command, args []string) {
	wallet, fromAddress, err := walletUnlock(cmd, fromFlag, passwordFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(fromAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}
	fund, ok := types.ParseCoinAmount(reserveFundInDTokenFlag)
	if !ok {
		utils.Error("Failed to parse fund")
	}
	col, ok := types.ParseCoinAmount(reserveCollateralInDTokenFlag)
	if !ok {
		utils.Error("Failed to parse collateral")
	}
	input := types.TxInput{
		Address: fromAddress,
		Coins: types.Coins{
			DneroWei: new(big.Int).SetUint64(0),
			DTokenWei: fund,
		},
		Sequence: uint64(seqFlag),
	}
	resourceIDs := []string{}
	for _, id := range resourceIDsFlag {
		resourceIDs = append(resourceIDs, id)
	}
	collateral := types.Coins{
		DneroWei: new(big.Int).SetUint64(0),
		DTokenWei: col,
	}
	if !collateral.IsPositive() {
		utils.Error("Invalid input: collateral must be positive\n")
	}

	reserveFundTx := &types.ReserveFundTx{
		Fee: types.Coins{
			DneroWei: new(big.Int).SetUint64(0),
			DTokenWei: fee,
		},
		Source:      input,
		ResourceIDs: resourceIDs,
		Collateral:  collateral,
		Duration:    durationFlag,
	}

	sig, err := wallet.Sign(fromAddress, reserveFundTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	reserveFundTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(reserveFundTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	var res *rpcc.RPCResponse
	if asyncFlag {
		res, err = client.Call("dnero.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	} else {
		res, err = client.Call("dnero.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	}
	if err != nil {
		utils.Error("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Server returned error: %v\n", res.Error)
	}
	fmt.Printf("Successfully broadcasted transaction.\n")
}

func init() {
	reserveFundCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	reserveFundCmd.Flags().StringVar(&fromFlag, "from", "", "Address to send from")
	reserveFundCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	reserveFundCmd.Flags().StringVar(&reserveFundInDTokenFlag, "fund", "0", "DToken amount to reserve")
	reserveFundCmd.Flags().StringVar(&reserveCollateralInDTokenFlag, "collateral", "0", "DToken amount as collateral")
	reserveFundCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeDTokenWeiNewFee), "Fee")
	reserveFundCmd.Flags().Uint64Var(&durationFlag, "duration", 1000, "Reserve duration")
	reserveFundCmd.Flags().StringSliceVar(&resourceIDsFlag, "resource_ids", []string{}, "Reserouce IDs")
	reserveFundCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")
	reserveFundCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")
	reserveFundCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")

	reserveFundCmd.MarkFlagRequired("chain")
	reserveFundCmd.MarkFlagRequired("from")
	reserveFundCmd.MarkFlagRequired("seq")
	reserveFundCmd.MarkFlagRequired("duration")
	reserveFundCmd.MarkFlagRequired("resource_id")
}

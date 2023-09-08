package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// scpCmd represents the scp command.
// Example:
//		dnerocli query scp --height=10
var scpCmd = &cobra.Command{
	Use:     "scp",
	Short:   "Get sentry candidate pool",
	Example: `dnerocli query scp --height=10`,
	Run:     doScpCmd,
}

func doScpCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	height := heightFlag
	res, err := client.Call("dnero.GetScpByHeight", rpc.GetScpByHeightArgs{Height: common.JSONUint64(height)})
	if err != nil {
		utils.Error("Failed to get sentry candidate pool: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get sentry candidate pool: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	scpCmd.Flags().Uint64Var(&heightFlag, "height", uint64(0), "height of the block")
	scpCmd.MarkFlagRequired("height")
}

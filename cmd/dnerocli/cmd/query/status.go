package query

import (
	"encoding/json"
	"fmt"

	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	"github.com/dnerochain/dnero/rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcc "github.com/ybbus/jsonrpc"
)

// statusCmd represents the account command.
// Example:
//		dnerocli query status
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Get blockchain status",
	Long:    `Get blockchain status.`,
	Example: `dnerocli query status`,
	Run: func(cmd *cobra.Command, args []string) {
		client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

		res, err := client.Call("dnero.GetStatus", rpc.GetStatusArgs{})
		if err != nil {
			utils.Error("Failed to get blockchain status: %v\n", err)
		}
		if res.Error != nil {
			utils.Error("Failed to retrieve blockchain status: %v\n", res.Error)
		}
		json, err := json.MarshalIndent(res.Result, "", "    ")
		if err != nil {
			utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
		}
		fmt.Println(string(json))
	},
}

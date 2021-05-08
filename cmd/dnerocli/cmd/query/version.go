package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	"github.com/dnerochain/dnero/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// versionCmd represents the version command.
// Example:
//		dnerocli query version
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get the Dnero version",
	Example: `dnerocli query version`,
	Run: func(cmd *cobra.Command, args []string) {
		client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

		res, err := client.Call("dnero.GetVersion", rpc.GetVersionArgs{})
		if err != nil {
			utils.Error("Failed to get version: %v\n", err)
		}
		if res.Error != nil {
			utils.Error("Failed to get version: %v\n", res.Error)
		}
		json, err := json.MarshalIndent(res.Result, "", "    ")
		if err != nil {
			utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
		}
		fmt.Println(string(json))
	},
}

package backup

import (
	"encoding/json"
	"fmt"

	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	"github.com/dnerochain/dnero/rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcc "github.com/ybbus/jsonrpc"
)

// snapshotCmd represents the snapshot backup command.
// Example:
//		dnerocli backup snapshot
var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Short:   "backup snapshot",
	Long:    `Backup snapshot.`,
	Example: `dnerocli backup snapshot`,
	Run:     doSnapshotCmd,
}

func doSnapshotCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("dnero.BackupSnapshot", rpc.BackupSnapshotArgs{Config: configFlag, Height: heightFlag})
	if err != nil {
		utils.Error("Failed to get backup snapshot call details: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get backup snapshot res details: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	snapshotCmd.Flags().StringVar(&configFlag, "config", "", "Config dir")
	snapshotCmd.MarkFlagRequired("config")
	snapshotCmd.Flags().Uint64Var(&heightFlag, "height", 0, "Snapshot height")
}

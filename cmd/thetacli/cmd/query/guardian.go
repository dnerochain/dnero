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

// guardianCmd retreves guardian related information from Dnero server.
// Example:
//		dnerocli query guardian
var guardianCmd = &cobra.Command{
	Use:     "guardian",
	Short:   "Get guardian info",
	Long:    `Get guardian status.`,
	Example: `dnerocli query guardian`,
	Run:     doGuardianCmd,
}

type GuardianResult struct {
	Address   string
	BlsPubkey string
	BlsPop    string
	Signature string
	Summary   string
}

func doGuardianCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("dnero.GetGuardianInfo", rpc.GetGuardianInfoArgs{})
	if err != nil {
		utils.Error("Failed to get guardian info: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get guardian info: %v\n", res.Error)
	}
	result := res.Result.(map[string]interface{})
	address, ok := result["Address"].(string)
	if !ok {
		json, err := json.MarshalIndent(res.Result, "", "    ")
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	blsPubkey, ok := result["BLSPubkey"].(string)
	if !ok {
		json, err := json.MarshalIndent(res.Result, "", "    ")
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	blsPop, ok := result["BLSPop"].(string)
	if !ok {
		json, err := json.MarshalIndent(res.Result, "", "    ")
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	sig, ok := result["Signature"].(string)
	if !ok {
		json, err := json.MarshalIndent(res.Result, "", "    ")
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	output := &GuardianResult{
		Address:   address,
		BlsPubkey: blsPubkey,
		BlsPop:    blsPop,
		Signature: sig,
	}
	output.Summary = address + blsPubkey + blsPop + sig
	json, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {}

package daemon

import (
	"context"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/dnerochain/dnero/cmd/dnerocli/rpc"
)

// startDaemonCmd runs the dnerocli daemon
// Example:
//		dnerocli daemon start --port=16889
var startDaemonCmd = &cobra.Command{
	Use:     "start",
	Short:   "Run the dnerocli daemon",
	Long:    `Run the dnerocli daemon.`,
	Example: `dnerocli daemon start --port=16889`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := cmd.Flag("config").Value.String()
		server, err := rpc.NewDneroCliRPCServer(cfgPath, portFlag)
		if err != nil {
			log.Fatalf("Failed to run the DneroCli Daemon: %v", err)
		}
		daemon := &DneroCliDaemon{
			RPC: server,
		}
		daemon.Start(context.Background())
		daemon.Wait()
	},
}

func init() {
	startDaemonCmd.Flags().StringVar(&portFlag, "port", "16889", "Port to run the DneroCli Daemon")
}

type DneroCliDaemon struct {
	RPC *rpc.DneroCliRPCServer

	// Life cycle
	wg      *sync.WaitGroup
	quit    chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

func (d *DneroCliDaemon) Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	d.ctx = c
	d.cancel = cancel

	if d.RPC != nil {
		d.RPC.Start(d.ctx)
	}
}

func (d *DneroCliDaemon) Stop() {
	d.cancel()
}

func (d *DneroCliDaemon) Wait() {
	if d.RPC != nil {
		d.RPC.Wait()
	}
}

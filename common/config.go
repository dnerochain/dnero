package common

import (
	"github.com/spf13/viper"
)

const (
	// CfgConfigPath defines custom config path
	CfgConfigPath = "config.path"

	// CfgDataPath defines custom DB path
	CfgDataPath = "data.path"

	// CfgKeyPath defines custom key path
	CfgKeyPath = "key.path"

	// CfgForceValidateSnapshot defines wether validation of snapshot can be skipped
	CfgForceValidateSnapshot = "snapshot.force_validate"

	// CfgGenesisHash defines the hash of the genesis block
	CfgGenesisHash = "genesis.hash"
	// CfgGenesisChainID defines the chainID.
	CfgGenesisChainID = "genesis.chainID"

	// CfgConsensusMaxEpochLength defines the maxium length of an epoch.
	CfgConsensusMaxEpochLength = "consensus.maxEpochLength"
	// CfgConsensusMinProposalWait defines the minimal interval between proposals.
	CfgConsensusMinProposalWait = "consensus.minProposalWait"
	// CfgConsensusMessageQueueSize defines the capacity of consensus message queue.
	CfgConsensusMessageQueueSize = "consensus.messageQueueSize"
	// CfgConsensusPassThroughGuardianVote defines the how guardian vote is handled.
	CfgConsensusPassThroughGuardianVote = "consensus.passThroughGuardianVote"

	// CfgStorageStatePruningEnabled indicates whether state pruning is enabled
	CfgStorageStatePruningEnabled = "storage.statePruningEnabled"
	// CfgStorageStatePruningInterval indicates the purning interval (in terms of blocks)
	CfgStorageStatePruningInterval = "storage.statePruningInterval"
	// CfgStorageStatePruningRetainedBlocks indicates the number of blocks prior to the latest finalized block to be retained
	CfgStorageStatePruningRetainedBlocks = "storage.statePruningRetainedBlocks"
	// CfgStorageStatePruningSkipCheckpoints indicates if the checkpoint state trie should be retained
	CfgStorageStatePruningSkipCheckpoints = "storage.statePruningSkipCheckpoints"
	// CfgStorageLevelDBCacheSize indicates Level DB cache size
	CfgStorageLevelDBCacheSize = "storage.levelDBCacheSize"
	// CfgStorageLevelDBHandles indicates Level DB handle count
	CfgStorageLevelDBHandles = "storage.levelDBHandles"

	// CfgSyncMessageQueueSize defines the capacity of Sync Manager message queue.
	CfgSyncMessageQueueSize = "sync.messageQueueSize"
	// CfgSyncDownloadByHash indicates whether should download blocks using hash.
	CfgSyncDownloadByHash = "sync.downloadByHash"
	// CfgSyncDownloadByHeader indicates whether should download blocks using header.
	CfgSyncDownloadByHeader = "sync.downloadByHeader"

	// CfgP2POpt sets which P2P network to use: p2p, libp2p, or both.
	CfgP2POpt = "p2p.opt"
	// CfgP2PReuseStream sets whether to reuse libp2p stream
	CfgP2PReuseStream = "p2p.reuseStream"
	// CfgP2PName sets the ID of local node in P2P network.
	CfgP2PName = "p2p.name"
	// CfgP2PVersion sets the version of P2P network.
	CfgP2PVersion = "p2p.version"
	// CfgP2PProtocolPrefix sets the protocol prefix of P2P network.
	CfgP2PProtocolPrefix = "p2p.protocolPrefix"
	// CfgP2PPort sets the port used by P2P network.
	CfgP2PPort = "p2p.port"
	// CfgP2PLPort sets the port used by P2P network.
	CfgP2PLPort = "p2p.libp2pPort"
	// CfgP2PSeeds sets the boostrap peers.
	CfgP2PSeeds = "p2p.seeds"
	// CfgLibP2PSeeds sets the boostrap peers in libp2p format.
	CfgLibP2PSeeds = "p2p.libp2pSeeds"
	// CfgLibP2PRendezvous is the libp2p rendezvous string
	CfgLibP2PRendezvous = "p2p.libp2pRendezvous"
	// CfgP2PMessageQueueSize sets the message queue size for network interface.
	CfgP2PMessageQueueSize = "p2p.messageQueueSize"
	// CfgP2PSeedPeerOnlyOutbound decides whether only the seed peers can be outbound peers.
	CfgP2PSeedPeerOnlyOutbound = "p2p.seedPeerOnlyOutbound"
	// CfgP2PSeedPeerOnly decides whether the node will connect to peers other than the seeds.
	CfgP2PSeedPeerOnly = "p2p.seedPeerOnly"
	// CfgP2PMinNumPeers specifies the minimal number of peers a node tries to maintain
	CfgP2PMinNumPeers = "p2p.minNumPeers"
	// CfgP2PMaxNumPeers specifies the maximal number of peers a node can simultaneously connected to
	CfgP2PMaxNumPeers = "p2p.maxNumPeers"
	// CfgMaxNumPersistentPeers sets the max number of peers to persist for normal nodes
	CfgMaxNumPersistentPeers = "p2p.maxNumPersistentPeers"
	// CfgP2PMaxNumPeersToBroadcast specifies the maximal number of peers to broadcast a message to
	CfgP2PMaxNumPeersToBroadcast = "p2p.maxNumPeersToBroadcast"
	// CfgBufferPoolSize defines the number of buffers in the pool.
	CfgBufferPoolSize = "p2p.bufferPoolSize"
	// CfgP2PConnectionFIFO specifies if the incoming connection policy is FIFO or LIFO
	CfgP2PConnectionFIFO = "p2p.connectionFIFO"
	// CfgP2PNatMapping sets whether to perform NAT mapping
	CfgP2PNatMapping = "p2p.natMapping"
	// CfgP2PMaxConnections specifies the number of max connections a node can accept
	CfgP2PMaxConnections = "p2p.maxConnections"

	// CfgSyncInboundResponseWhitelist filters inbound messages based on peer ID.
	CfgSyncInboundResponseWhitelist = "sync.inboundResponseWhitelist"

	// CfgRPCEnabled sets whether to run RPC service.
	CfgRPCEnabled = "rpc.enabled"
	// CfgRPCAddress sets the binding address of RPC service.
	CfgRPCAddress = "rpc.address"
	// CfgRPCPort sets the port of RPC service.
	CfgRPCPort = "rpc.port"
	// CfgRPCMaxConnections limits concurrent connections accepted by RPC server.
	CfgRPCMaxConnections = "rpc.maxConnections"
	// CfgRPCTimeoutSecs set a timeout for RPC.
	CfgRPCTimeoutSecs = "rpc.timeoutSecs"

	// CfgLogLevels sets the log level.
	CfgLogLevels = "log.levels"
	// CfgLogPrintSelfID determines whether to print node's ID in log (Useful in simulation when
	// there are more than one node running).
	CfgLogPrintSelfID = "log.printSelfID"

	// CfgGuardianRoundLength defines the length of a guardian voting round.
	CfgGuardianRoundLength = "guardian.roundLength"

	// Graphite Server to collet metrics
	CfgMetricsServer = "metrics.server"

	// CfgProfEnabled to enable profiling
	CfgProfEnabled = "prof.enabled"

	// CfgForceGCEnabled to enable force GC
	CfgForceGCEnabled = "gc.enabled"
)

// Starting block heights of features.
const (
	FeatureGuardian uint64 = 0
)

// InitialConfig is the default configuration produced by init command.
const InitialConfig = `# Dnero configuration
//p2p: //DEL
  //port: 5000
  //seeds: 127.0.0.1:6000,127.0.0.1:7000
p2p:
  port: 6000
  seeds: 127.0.0.1:7000,127.0.0.1:8000
`

func init() {
	viper.SetDefault(CfgForceValidateSnapshot, false)

	viper.SetDefault(CfgConsensusMaxEpochLength, 20)
	viper.SetDefault(CfgConsensusMinProposalWait, 6)
	viper.SetDefault(CfgConsensusMessageQueueSize, 512)
	viper.SetDefault(CfgConsensusPassThroughGuardianVote, false)

	viper.SetDefault(CfgSyncMessageQueueSize, 512)
	viper.SetDefault(CfgSyncDownloadByHash, false)
	viper.SetDefault(CfgSyncDownloadByHeader, true)

	viper.SetDefault(CfgStorageStatePruningEnabled, true)
	viper.SetDefault(CfgStorageStatePruningInterval, 16)
	viper.SetDefault(CfgStorageStatePruningRetainedBlocks, 2048)
	viper.SetDefault(CfgStorageStatePruningSkipCheckpoints, true)
	viper.SetDefault(CfgStorageLevelDBCacheSize, 256)
	viper.SetDefault(CfgStorageLevelDBHandles, 16)

	viper.SetDefault(CfgRPCEnabled, false)
	viper.SetDefault(CfgP2PMessageQueueSize, 512)
	viper.SetDefault(CfgP2PName, "Anonymous")
	//viper.SetDefault(CfgP2PPort, 50001)//DEL
	viper.SetDefault(CfgP2PPort, 61000)
	viper.SetDefault(CfgP2PSeeds, "") //TODO: Hardcode main seeds
	viper.SetDefault(CfgP2PSeedPeerOnlyOutbound, false)
	//viper.SetDefault(CfgP2POpt, P2POptLibp2p) // FIXME: this for some reason doesn't work
	viper.SetDefault(CfgP2POpt, 0)
	viper.SetDefault(CfgP2PReuseStream, true)
	viper.SetDefault(CfgP2PSeedPeerOnly, false)
	viper.SetDefault(CfgP2PMinNumPeers, 32)
	//viper.SetDefault(CfgP2PMaxNumPeers, 256)
	viper.SetDefault(CfgP2PMaxNumPeers, 64)
	viper.SetDefault(CfgP2PMaxNumPeersToBroadcast, 64)
	viper.SetDefault(CfgMaxNumPersistentPeers, 10)
	viper.SetDefault(CfgBufferPoolSize, 8)
	viper.SetDefault(CfgP2PConnectionFIFO, false)
	viper.SetDefault(CfgP2PNatMapping, false)
	viper.SetDefault(CfgP2PMaxConnections, 2048)

	viper.SetDefault(CfgRPCAddress, "0.0.0.0")
	//viper.SetDefault(CfgRPCPort, "16888")//DEL
	viper.SetDefault(CfgRPCPort, "15511")
	viper.SetDefault(CfgRPCMaxConnections, 200)
	viper.SetDefault(CfgRPCTimeoutSecs, 60)

	viper.SetDefault(CfgLogLevels, "*:debug")
	viper.SetDefault(CfgLogPrintSelfID, false)

	viper.SetDefault(CfgGuardianRoundLength, 30)

	//TODO: Setup Dnero own Metrics Server
	//viper.SetDefault(CfgMetricsServer, "guardian-metrics.dnerochain.org")

	viper.SetDefault(CfgProfEnabled, false)
	viper.SetDefault(CfgForceGCEnabled, true)
}

// WriteInitialConfig writes initial config file to file system.
func WriteInitialConfig(filePath string) error {
	return WriteFileAtomic(filePath, []byte(InitialConfig), 0600)
}

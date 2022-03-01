package common

// HeightEnableValidatorReward specifies the minimal block height to enable the validator DTOKEN reward
const HeightEnableValidatorReward uint64 = 10001 // block #10001

// HeightEnableDneroV1 specifies the minimal block height to enable the DneroV1.0 feature.
const HeightEnableDneroV1 uint64 = 20001 // block #20001

//DEL
// HeightLowerGNStakeThresholdTo1000 specifies the minimal block height to lower the GN Stake Threshold to 1,000 DNERO
//const HeightLowerGNStakeThresholdTo1000 uint64 = ### // StakeDeposit Fork Removed
//DEL

// HeightEnableSmartContract specifies the minimal block height to enable the Turing-complete smart contract support
const HeightEnableSmartContract uint64 = 50001 // block #50001

// HeightSampleStakingReward specifies the block heigth to enable sampling of staking reward
const HeightSampleStakingReward uint64 = 100001 // block #100001

// HeightNewFeeAdjustment specifies the block heigth to enable transaction fee burning adjustment
const HeightNewFeeAdjustment uint64 = 120001 // block #120001

// HeightEnableDneroV2 specifies the minimal block height to enable the DneroV2.0 feature.
const HeightEnableDneroV2 uint64 = 200001 // block #200001

// HeightRPCCompatibility specifies the block height to enable Ethereum compatible RPC support
const HeightRPCCompatibility uint64 = 210001 // block #210001

// HeightTxWrapperExtension specifies the block height to extend the Tx Wrapper
const HeightTxWrapperExtension uint64 = 220001 // block #220001

// HeightSupportDneroTokenInSmartContract specifies the block height to support Dnero in smart contracts
const HeightSupportDneroTokenInSmartContract uint64 = 250001 // block #250001

// CheckpointInterval defines the interval between checkpoints.
const CheckpointInterval = int64(100)

// IsCheckPointHeight returns if a block height is a checkpoint.
func IsCheckPointHeight(height uint64) bool {
	return height%uint64(CheckpointInterval) == 1
}

// LastCheckPointHeight returns the height of the last checkpoint
func LastCheckPointHeight(height uint64) uint64 {
	multiple := height / uint64(CheckpointInterval)
	lastCheckpointHeight := uint64(CheckpointInterval)*multiple + 1
	return lastCheckpointHeight
}

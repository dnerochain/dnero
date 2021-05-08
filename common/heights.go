package common

// HeightEnableValidatorReward specifies the minimal block height to enable the validtor DFUEL reward
const HeightEnableValidatorReward uint64 = 101 // block #101

// HeightEnableDneroV1 specifies the minimal block height to enable the DneroV1 feature.
const HeightEnableDneroV1 uint64 = 102 // block #102

// HeightLowerGNStakeThresholdTo1000 specifies the minimal block height to lower the GN Stake Threshold to 1,000 DNERO
const HeightLowerGNStakeThresholdTo1000 uint64 = 103 // block #103

// HeightEnableSmartContract specifies the minimal block height to eanble the Turing-complete smart contract support
const HeightEnableSmartContract uint64 = 104 // block #104

// HeightSampleStakingReward specifies the block heigth to enable sampling of staking reward
const HeightSampleStakingReward uint64 = 105 // block #105

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

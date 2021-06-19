package common

// HeightEnableDneroV1 specifies the minimal block height to enable the DneroV1 feature.
const HeightEnableDneroV1 uint64 = 101 // TODO MAINNET: start HeightEnableDneroV1 at block #101

// HeightEnableValidatorReward specifies the minimal block height to enable the Validator DFUEL reward
const HeightEnableValidatorReward uint64 = 201 // TODO MAINNET: start HeightEnableValidatorReward at block #201

// HeightLowerGNStakeThresholdTo100 specifies the minimal block height to lower the GN Stake Threshold to 100 DNERO
//const HeightLowerGNStakeThresholdTo100 uint64 = ### // block #000 //GN StakeDeposit Fork Removed

// HeightEnableSmartContract specifies the minimal block height to enable the Turing-complete smart contract support
const HeightEnableSmartContract uint64 = 401 // block #401

// HeightSampleStakingReward specifies the block height to enable sampling of staking reward
const HeightSampleStakingReward uint64 = 501 // block #501

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

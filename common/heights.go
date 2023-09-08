package common


// HeightEnableDneroV1 specifies the minimal block height to enable the DneroV1.0 feature.
const HeightEnableDneroV1 uint64 = 1001 // block #1001

// HeightEnableValidatorReward specifies the minimal block height to enable the validtor DTOKEN reward
const HeightEnableValidatorReward uint64 = 50001 // block #50001

//DEL
// HeightLowerGNStakeThresholdTo1000 specifies the minimal block height to lower the GN Stake Threshold to 1,000 DNERO
//const HeightLowerGNStakeThresholdTo1000 uint64 = ### // StakeDeposit Fork Removed
//DEL

// HeightEnableSmartContract specifies the minimal block height to eanble the Turing-complete smart contract support
const HeightEnableSmartContract uint64 = 60001 // block #60001

// HeightRPCCompatibility specifies the block height to enable Ethereum compatible RPC support
const HeightRPCCompatibility uint64 = 70001 // block #70001

// HeightNewFeeAdjustment specifies the block heigth to enable transaction fee burning adjustment
const HeightNewFeeAdjustment uint64 = 80001 // block #80001

// HeightSampleStakingReward specifies the block heigth to enable sampling of staking reward
const HeightSampleStakingReward uint64 = 90001 // block #90001

// HeightEnableDneroV2 specifies the minimal block height to enable the DneroV2.0 feature.
const HeightEnableDneroV2 uint64 = 100001 // block #100001

// HeightTxWrapperExtension specifies the block height to extend the Tx Wrapper
const HeightTxWrapperExtension uint64 = 110001 // block #110001

// HeightSupportDneroTokenInSmartContract specifies the block height to support Dnero in smart contracts
const HeightSupportDneroTokenInSmartContract uint64 = 120001 // block #120001

//DEL
// HeightValidatorStakeChangedTo200K specifies the block height to lower the validator stake to 200,000 Dnero
//const HeightValidatorStakeChangedTo200K uint64 = ### // ValidatorStake Fork Removed
//DEL

// HeightSupportWrappedDnero specifies the block height to support wrapped Dnero
const HeightSupportWrappedDnero uint64 = 150001 // block #150001

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

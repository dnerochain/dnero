package types

import (
	"math/big"

	"github.com/dnerochain/dnero/common"
)

const (
	// DenomDneroWei is the basic unit of dnero, 1 Dnero = 10^18 DneroWei
	DenomDneroWei string = "DneroWei"

	// DenomDTokenWei is the basic unit of dnero, 1 Dnero = 10^18 DneroWei
	DenomDTokenWei string = "DTokenWei"

	// Initial gas parameters

	// MinimumGasPrice is the minimum gas price for a smart contract transaction
	MinimumGasPrice uint64 = 1e8

	// MaximumTxGasLimit is the maximum gas limit for a smart contract transaction
	//MaximumTxGasLimit uint64 = 2e6
	MaximumTxGasLimit uint64 = 10e6

	// MinimumTransactionFeeDTokenWei specifies the minimum fee for a regular transaction
	MinimumTransactionFeeDTokenWei uint64 = 1e12

	// NewFee gas burn adjustment

	// MinimumGasPrice is the minimum gas price for a smart contract transaction
	MinimumGasPriceNewFee uint64 = 4e12

	// MaximumTxGasLimit is the maximum gas limit for a smart contract transaction
	MaximumTxGasLimitNewFee uint64 = 20e6

	// MinimumTransactionFeeDTokenWei specifies the minimum fee for a regular transaction
	MinimumTransactionFeeDTokenWeiNewFee uint64 = 3e17

	// MaxAccountsAffectedPerTx specifies the max number of accounts one transaction is allowed to modify to avoid spamming
	MaxAccountsAffectedPerTx = 512
)

const (
	// ValidatorDneroGenerationRateNumerator is used for calculating the generation rate of Dnero for validators
	//ValidatorDneroGenerationRateNumerator int64 = 317
	ValidatorDneroGenerationRateNumerator int64 = 0 // ZERO inflation for Dnero

	// ValidatorDneroGenerationRateDenominator is used for calculating the generation rate of Dnero for validators
	// ValidatorDneroGenerationRateNumerator / ValidatorDneroGenerationRateDenominator is the amount of DneroWei
	// generated per existing DneroWei per new block
	ValidatorDneroGenerationRateDenominator int64 = 1e11

	// ValidatorDTokenGenerationRateNumerator is used for calculating the generation rate of DToken for validators
	ValidatorDTokenGenerationRateNumerator int64 = 0 // ZERO initial inflation for DToken

	// ValidatorDTokenGenerationRateDenominator is used for calculating the generation rate of DToken for validators
	// ValidatorDTokenGenerationRateNumerator / ValidatorDTokenGenerationRateDenominator is the amount of DTokenWei
	// generated per existing DneroWei per new block
	ValidatorDTokenGenerationRateDenominator int64 = 1e9

	// RegularDTokenGenerationRateNumerator is used for calculating the generation rate of DToken for other types of accounts
	//RegularDTokenGenerationRateNumerator int64 = 1900
	RegularDTokenGenerationRateNumerator int64 = 0 // ZERO initial inflation for DToken

	// RegularDTokenGenerationRateDenominator is used for calculating the generation rate of DToken for other types of accounts
	// RegularDTokenGenerationRateNumerator / RegularDTokenGenerationRateDenominator is the amount of DTokenWei
	// generated per existing DneroWei per new block
	RegularDTokenGenerationRateDenominator int64 = 1e10
)

const (

	// ServiceRewardVerificationBlockDelay gives the block delay for service certificate verification
	ServiceRewardVerificationBlockDelay uint64 = 2

	// ServiceRewardFulfillmentBlockDelay gives the block delay for service reward fulfillment
	ServiceRewardFulfillmentBlockDelay uint64 = 4
)

const (

	// MaximumTargetAddressesForStakeBinding gives the maximum number of target addresses that can be associated with a bound stake
	MaximumTargetAddressesForStakeBinding uint = 1024

	// MaximumFundReserveDuration indicates the maximum duration (in terms of number of blocks) of reserving fund
	MaximumFundReserveDuration uint64 = 12 * 3600

	// MinimumFundReserveDuration indicates the minimum duration (in terms of number of blocks) of reserving fund
	MinimumFundReserveDuration uint64 = 300

	// ReservedFundFreezePeriodDuration indicates the freeze duration (in terms of number of blocks) of the reserved fund
	ReservedFundFreezePeriodDuration uint64 = 5
)

func GetMinimumGasPrice(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightNewFeeAdjustment {
		return new(big.Int).SetUint64(MinimumGasPrice)
	}

	return new(big.Int).SetUint64(MinimumGasPriceNewFee)
}

func GetMaxGasLimit(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightNewFeeAdjustment {
		return new(big.Int).SetUint64(MaximumTxGasLimit)
	}

	return new(big.Int).SetUint64(MaximumTxGasLimitNewFee)
}

func GetMinimumTransactionFeeDTokenWei(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightNewFeeAdjustment {
		return new(big.Int).SetUint64(MinimumTransactionFeeDTokenWei)
	}

	return new(big.Int).SetUint64(MinimumTransactionFeeDTokenWeiNewFee)
}

// Special handling for many-to-many SendTx
func GetSendTxMinimumTransactionFeeDTokenWei(numAccountsAffected uint64, blockHeight uint64) *big.Int {
	if blockHeight < common.HeightNewFeeAdjustment {
		return new(big.Int).SetUint64(MinimumTransactionFeeDTokenWei) // backward compatiblity
	}

	if numAccountsAffected < 2 {
		numAccountsAffected = 2
	}

	// minSendTxFee = numAccountsAffected * MinimumTransactionFeeDTokenWeiNewFee / 2
	minSendTxFee := big.NewInt(1).Mul(new(big.Int).SetUint64(numAccountsAffected), new(big.Int).SetUint64(MinimumTransactionFeeDTokenWeiNewFee))
	minSendTxFee = big.NewInt(1).Div(minSendTxFee, new(big.Int).SetUint64(2))

	return minSendTxFee
}

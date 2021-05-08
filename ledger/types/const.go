package types

const (
	// DenomDneroWei is the basic unit of dnero, 1 Dnero = 10^18 DneroWei
	DenomDneroWei string = "DneroWei"

	// DenomDFuelWei is the basic unit of dnero, 1 Dnero = 10^18 DneroWei
	DenomDFuelWei string = "DFuelWei"

	// MinimumGasPrice is the minimum gas price for a smart contract transaction
	MinimumGasPrice uint64 = 1e8

	// MaximumTxGasLimit is the maximum gas limit for a smart contract transaction
	//MaximumTxGasLimit uint64 = 2e6
	MaximumTxGasLimit uint64 = 10e6

	// MinimumTransactionFeeDFuelWei specifies the minimum fee for a regular transaction
	MinimumTransactionFeeDFuelWei uint64 = 1e12

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

	// ValidatorDFuelGenerationRateNumerator is used for calculating the generation rate of DFuel for validators
	ValidatorDFuelGenerationRateNumerator int64 = 0 // ZERO initial inflation for DFuel

	// ValidatorDFuelGenerationRateDenominator is used for calculating the generation rate of DFuel for validators
	// ValidatorDFuelGenerationRateNumerator / ValidatorDFuelGenerationRateDenominator is the amount of DFuelWei
	// generated per existing DneroWei per new block
	ValidatorDFuelGenerationRateDenominator int64 = 1e9

	// RegularDFuelGenerationRateNumerator is used for calculating the generation rate of DFuel for other types of accounts
	//RegularDFuelGenerationRateNumerator int64 = 1900
	RegularDFuelGenerationRateNumerator int64 = 0 // ZERO initial inflation for DFuel

	// RegularDFuelGenerationRateDenominator is used for calculating the generation rate of DFuel for other types of accounts
	// RegularDFuelGenerationRateNumerator / RegularDFuelGenerationRateDenominator is the amount of DFuelWei
	// generated per existing DneroWei per new block
	RegularDFuelGenerationRateDenominator int64 = 1e10
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

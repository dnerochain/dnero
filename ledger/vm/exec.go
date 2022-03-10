package vm

import (
	"math"
	"math/big"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/core"
	"github.com/dnerochain/dnero/crypto"
	"github.com/dnerochain/dnero/ledger/state"
	"github.com/dnerochain/dnero/ledger/types"
	"github.com/dnerochain/dnero/ledger/vm/params"
)

// Execute executes the given smart contract
func Execute(parentBlock *core.Block, tx *types.SmartContractTx, storeView *state.StoreView) (evmRet common.Bytes,
	contractAddr common.Address, gasUsed uint64, evmErr error) {
	context := Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		Origin:      tx.From.Address,
		GasPrice:    tx.GasPrice,
		GasLimit:    tx.GasLimit,
		BlockNumber: new(big.Int).SetUint64(parentBlock.Height + 1),
		Time:        parentBlock.Timestamp,
		Difficulty:  new(big.Int).SetInt64(0),
	}
	chainIDBigInt := mapChainID(parentBlock.ChainID)
	chainConfig := &params.ChainConfig{
		ChainID: chainIDBigInt,
	}
	config := Config{}
	evm := NewEVM(context, storeView, chainConfig, config)

	value := tx.From.Coins.DTokenWei
	if value == nil {
		value = big.NewInt(0)
	}
	gasLimit := tx.GasLimit
	fromAddr := tx.From.Address
	contractAddr = tx.To.Address
	createContract := (contractAddr == common.Address{})

	// if gasLimit > maxGasLimit {
	// 	return common.Bytes{}, common.Address{}, 0, ErrInvalidGasLimit
	// }
	blockHeight := storeView.Height() + 1
	maxGasLimit := types.GetMaxGasLimit(blockHeight)
	if new(big.Int).SetUint64(gasLimit).Cmp(maxGasLimit) > 0 {
		return common.Bytes{}, common.Address{}, 0, ErrInvalidGasLimit
	}

	intrinsicGas, err := calculateIntrinsicGas(tx.Data, createContract)
	if err != nil {
		return common.Bytes{}, common.Address{}, 0, err
	}
	if intrinsicGas > gasLimit {
		return common.Bytes{}, common.Address{}, 0, ErrOutOfGas
	}

	var leftOverGas uint64
	remainingGas := gasLimit - intrinsicGas
	if createContract {
		code := tx.Data
		evmRet, contractAddr, leftOverGas, evmErr = evm.Create(AccountRef(fromAddr), code, remainingGas, value)
	} else {
		input := tx.Data
		evmRet, leftOverGas, evmErr = evm.Call(AccountRef(fromAddr), contractAddr, input, remainingGas, value)
	}

	if leftOverGas > gasLimit { // should not happen
		gasUsed = uint64(0)
	} else {
		gasUsed = gasLimit - leftOverGas
	}

	return evmRet, contractAddr, gasUsed, evmErr
}

// calculateIntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func calculateIntrinsicGas(data []byte, createContract bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if createContract {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGas < nz {
			return 0, ErrOutOfGas
		}
		gas += nz * params.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrOutOfGas
		}
		gas += z * params.TxDataZeroGas
	}
	return gas, nil
}

// To be compatible with Ethereum, mapChainID() returns 1 for "mainnet", 3 for "testnet_sapphire", and 4 for "testnet_amber"
// Reference: https://github.com/ethereum/go-ethereum/blob/43cd31ea9f57e26f8f67aa8bd03bbb0a50814465/params/config.go#L55
func mapChainID(chainIDStr string) *big.Int {
	if chainIDStr == "mainnet" { // correspond to the Ethereum mainnet
		return big.NewInt(1)
	} else if chainIDStr == "testnet_sapphire" { // correspond to Ropsten
		return big.NewInt(3)
	} else if chainIDStr == "testnet_amber" { // correspond to Rinkeby
		return big.NewInt(4)
	} else if chainIDStr == "testnet" {
		return big.NewInt(5)
	} else if chainIDStr == "privatenet" {
		return big.NewInt(6)
	}

	chainIDBigInt := new(big.Int).Abs(crypto.Keccak256Hash(common.Bytes(chainIDStr)).Big()) // all other chainIDs
	return chainIDBigInt
}

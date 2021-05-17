package execution

import (
	"fmt"
	"math/big"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/common/result"
	"github.com/dnerochain/dnero/core"
	st "github.com/dnerochain/dnero/ledger/state"
	"github.com/dnerochain/dnero/ledger/types"
)

var _ TxExecutor = (*DepositStakeExecutor)(nil)

// ------------------------------- DepositStake Transaction -----------------------------------

// DepositStakeExecutor implements the TxExecutor interface
type DepositStakeExecutor struct {
}

// NewDepositStakeExecutor creates a new instance of DepositStakeExecutor
func NewDepositStakeExecutor() *DepositStakeExecutor {
	return &DepositStakeExecutor{}
}

func (exec *DepositStakeExecutor) sanityCheck(chainID string, view *st.StoreView, transaction types.Tx) result.Result {
	// Feature block height check
	blockHeight := view.Height() + 1 // the view points to the parent of the current block
	if _, ok := transaction.(*types.DepositStakeTxV1); ok && blockHeight < common.HeightEnableDneroV1 {
		return result.Error("Feature guardian is not active yet")
	}

	tx := exec.castTx(transaction)

	res := tx.Source.ValidateBasic()
	if res.IsError() {
		return res
	}

	sourceAccount, success := getInput(view, tx.Source)
	if success.IsError() {
		return result.Error("Failed to get the source account: %v", tx.Source.Address)
	}

	signBytes := tx.SignBytes(chainID)
	res = validateInputAdvanced(sourceAccount, signBytes, tx.Source)
	if res.IsError() {
		logger.Debugf(fmt.Sprintf("validateSourceAdvanced failed on %v: %v", tx.Source.Address.Hex(), res))
		return res
	}

	if !sanityCheckForFee(tx.Fee) {
		return result.Error("Insufficient fee. Transaction fee needs to be at least %v DFuelWei",
			types.MinimumTransactionFeeDFuelWei).WithErrorCode(result.CodeInvalidFee)
	}

	if !(tx.Purpose == core.StakeForValidator || tx.Purpose == core.StakeForGuardian) {
		return result.Error("Invalid stake purpose!").
			WithErrorCode(result.CodeInvalidStakePurpose)
	}

	stake := tx.Source.Coins.NoNil()
	if !stake.IsValid() || !stake.IsNonnegative() {
		return result.Error("Invalid stake for stake deposit!").
			WithErrorCode(result.CodeInvalidStake)
	}

	if stake.DFuelWei.Cmp(types.Zero) != 0 {
		return result.Error("DFuel has to be zero for stake deposit!").
			WithErrorCode(result.CodeInvalidStake)
	}

	// Minimum stake deposit requirement to avoid spamming
	if tx.Purpose == core.StakeForValidator && stake.DneroWei.Cmp(core.MinValidatorStakeDeposit) < 0 {
		return result.Error("Insufficient amount of stake, at least %v DneroWei is required for each validator deposit", core.MinValidatorStakeDeposit).
			WithErrorCode(result.CodeInsufficientStake)
	}

	if tx.Purpose == core.StakeForGuardian {
		minGuardianStake := core.MinGuardianStakeDeposit
		//if blockHeight >= common.HeightLowerGNStakeThresholdTo100 { //StakeDeposit Fork Removed
			//minGuardianStake = core.MinGuardianStakeDeposit100
		//}
		if stake.DneroWei.Cmp(minGuardianStake) < 0 {
			return result.Error("Insufficient amount of stake, at least %v DneroWei is required for each guardian deposit", minGuardianStake).
				WithErrorCode(result.CodeInsufficientStake)
		}
	}

	minimalBalance := stake.Plus(tx.Fee)
	if !sourceAccount.Balance.IsGTE(minimalBalance) {
		logger.Infof(fmt.Sprintf("DepositStake: Source did not have enough balance %v", tx.Source.Address.Hex()))
		return result.Error("DepositStake: Source balance is %v, but required minimal balance is %v",
			sourceAccount.Balance, minimalBalance).WithErrorCode(result.CodeInsufficientStake)
	}

	return result.OK
}

func (exec *DepositStakeExecutor) process(chainID string, view *st.StoreView, transaction types.Tx) (common.Hash, result.Result) {
	blockHeight := view.Height() + 1 // the view points to the parent of the current block

	tx := exec.castTx(transaction)

	sourceAccount, success := getInput(view, tx.Source)
	if success.IsError() {
		return common.Hash{}, result.Error("Failed to get the source account")
	}

	if !chargeFee(sourceAccount, tx.Fee) {
		return common.Hash{}, result.Error("Failed to charge transaction fee")
	}

	stake := tx.Source.Coins.NoNil()
	if !sourceAccount.Balance.IsGTE(stake) {
		return common.Hash{}, result.Error("Not enough balance to stake").WithErrorCode(result.CodeNotEnoughBalanceToStake)
	}

	sourceAddress := tx.Source.Address
	holderAddress := tx.Holder.Address

	if tx.Purpose == core.StakeForValidator {
		sourceAccount.Balance = sourceAccount.Balance.Minus(stake)
		stakeAmount := stake.DneroWei
		vcp := view.GetValidatorCandidatePool()
		err := vcp.DepositStake(sourceAddress, holderAddress, stakeAmount)
		if err != nil {
			return common.Hash{}, result.Error("Failed to deposit stake, err: %v", err)
		}
		view.UpdateValidatorCandidatePool(vcp)
	} else if tx.Purpose == core.StakeForGuardian {
		sourceAccount.Balance = sourceAccount.Balance.Minus(stake)
		stakeAmount := stake.DneroWei
		gcp := view.GetGuardianCandidatePool()

		if !gcp.Contains(holderAddress) {
			if tx.BlsPubkey.IsEmpty() {
				return common.Hash{}, result.Error("Must provide BLS Pubkey")
			}
			if tx.BlsPop.IsEmpty() {
				return common.Hash{}, result.Error("Must provide BLS POP")
			}
			if tx.HolderSig == nil || tx.HolderSig.IsEmpty() {
				return common.Hash{}, result.Error("Must provide Holder Signature")
			}

			if !tx.HolderSig.Verify(tx.BlsPop.ToBytes(), tx.Holder.Address) {
				return common.Hash{}, result.Error("BLS key info is not properly signed")
			}

			if !tx.BlsPop.PopVerify(tx.BlsPubkey) {
				return common.Hash{}, result.Error("BLS pop is invalid")
			}
		}

		err := gcp.DepositStake(sourceAddress, holderAddress, stakeAmount, tx.BlsPubkey, blockHeight)
		if err != nil {
			return common.Hash{}, result.Error("Failed to deposit stake, err: %v", err)
		}
		view.UpdateGuardianCandidatePool(gcp)
	} else {
		return common.Hash{}, result.Error("Invalid staking purpose").WithErrorCode(result.CodeInvalidStakePurpose)
	}

	// Only update stake transaction height list for validator stake tx.
	if tx.Purpose == core.StakeForValidator {
		hl := view.GetStakeTransactionHeightList()
		if hl == nil {
			hl = &types.HeightList{}
		}
		blockHeight := view.Height() + 1 // the view points to the parent of the current block
		hl.Append(blockHeight)
		view.UpdateStakeTransactionHeightList(hl)
	}

	sourceAccount.Sequence++
	view.SetAccount(sourceAddress, sourceAccount)

	txHash := types.TxID(chainID, tx)
	return txHash, result.OK
}

func (exec *DepositStakeExecutor) getTxInfo(transaction types.Tx) *core.TxInfo {
	tx := exec.castTx(transaction)
	return &core.TxInfo{
		Address:           tx.Source.Address,
		Sequence:          tx.Source.Sequence,
		EffectiveGasPrice: exec.calculateEffectiveGasPrice(transaction),
	}
}

func (exec *DepositStakeExecutor) calculateEffectiveGasPrice(transaction types.Tx) *big.Int {
	tx := exec.castTx(transaction)
	fee := tx.Fee
	gas := new(big.Int).SetUint64(types.GasDepositStakeTx)
	effectiveGasPrice := new(big.Int).Div(fee.DFuelWei, gas)
	return effectiveGasPrice
}

func (exec *DepositStakeExecutor) castTx(transaction types.Tx) *types.DepositStakeTxV1 {
	if tx, ok := transaction.(*types.DepositStakeTxV1); ok {
		return tx
	}
	if tx, ok := transaction.(*types.DepositStakeTx); ok {
		return &types.DepositStakeTxV1{
			Fee:     tx.Fee,
			Source:  tx.Source,
			Holder:  tx.Holder,
			Purpose: tx.Purpose,
		}
	}
	panic("Unreachable code")
}

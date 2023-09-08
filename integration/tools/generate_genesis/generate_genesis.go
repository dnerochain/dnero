package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/core"
	"github.com/dnerochain/dnero/ledger/state"
	"github.com/dnerochain/dnero/ledger/types"
	"github.com/dnerochain/dnero/rlp"
	"github.com/dnerochain/dnero/store/database/backend"
	"github.com/dnerochain/dnero/store/trie"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "genesis"})

const (
	GenBlockHashMode int = iota
	GenGenesisFileMode
)

type StakeDeposit struct {
	Source string `json:"source"`
	Holder string `json:"holder"`
	Amount string `json:"amount"`
}

//
// Example:
// pushd $DNERO_HOME/integration/privatenet/node
// generate_genesis -chainID=privatenet -allocationsnapshot=./data/genesis_dnero_allocation_snapshot.json -stake_deposit=./data/genesis_stake_deposit.json -genesis=./genesis
//
func main() {
	chainID, allocationSnapshotJSONFilePath, stakeDepositFilePath, genesisSnapshotFilePath := parseArguments()

	sv, metadata, err := generateGenesisSnapshot(chainID, allocationSnapshotJSONFilePath, stakeDepositFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate genesis snapshot: %v", err))
	}

	err = sanityChecks(sv)
	if err != nil {
		panic(fmt.Sprintf("Sanity checks failed: %v", err))
	} else {
		logger.Infof("Sanity checks all passed.")
	}

	err = writeGenesisSnapshot(sv, metadata, genesisSnapshotFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to write genesis snapshot: %v", err))
	}

	genesisBlockHeader := metadata.TailTrio.Second.Header
	genesisBlockHash := genesisBlockHeader.Hash()

	fmt.Println("")
	fmt.Printf("--------------------------------------------------------------------------\n")
	fmt.Printf("Genesis block hash: %v\n", genesisBlockHash.Hex())
	fmt.Printf("--------------------------------------------------------------------------\n")
	fmt.Println("")
}

func parseArguments() (chainID, allocationSnapshotJSONFilePath, stakeDepositFilePath, genesisSnapshotFilePath string) {
	chainIDPtr := flag.String("chainID", "local_chain", "the ID of the chain")
	allocationSnapshotJSONFilePathPtr := flag.String("allocationsnapshot", "./dnero_allocation_snapshot.json", "the json file contain the ALLOCATION balance snapshot")
	stakeDepositFilePathPtr := flag.String("stake_deposit", "./stake_deposit.json", "the initial stake deposits")
	genesisSnapshotFilePathPtr := flag.String("genesis", "./genesis", "the genesis snapshot")
	flag.Parse()

	chainID = *chainIDPtr
	allocationSnapshotJSONFilePath = *allocationSnapshotJSONFilePathPtr
	stakeDepositFilePath = *stakeDepositFilePathPtr
	genesisSnapshotFilePath = *genesisSnapshotFilePathPtr

	return
}

// generateGenesisSnapshot generates the genesis snapshot.
func generateGenesisSnapshot(chainID, allocationSnapshotJSONFilePath, stakeDepositFilePath string) (*state.StoreView, *core.SnapshotMetadata, error) {
	metadata := &core.SnapshotMetadata{}
	genesisHeight := core.GenesisBlockHeight

	sv := loadInitialBalances(allocationSnapshotJSONFilePath)
	performInitialStakeDeposit(stakeDepositFilePath, genesisHeight, sv)

	stateHash := sv.Hash()

	genesisBlock := core.NewBlock()
	genesisBlock.ChainID = chainID
	genesisBlock.Height = genesisHeight
	genesisBlock.Epoch = genesisBlock.Height
	genesisBlock.Parent = common.Hash{}
	genesisBlock.StateHash = stateHash
	genesisBlock.Timestamp = big.NewInt(time.Now().Unix())

	metadata.TailTrio = core.SnapshotBlockTrio{
		First:  core.SnapshotFirstBlock{},
		Second: core.SnapshotSecondBlock{Header: genesisBlock.BlockHeader},
		Third:  core.SnapshotThirdBlock{},
	}

	return sv, metadata, nil
}

func loadInitialBalances(allocationSnapshotJSONFilePath string) *state.StoreView {
	initDTokenToDneroRatio := new(big.Int).SetUint64(5)
	sv := state.NewStoreView(0, common.Hash{}, backend.NewMemDatabase())

	allocationSnapshotJSONFile, err := os.Open(allocationSnapshotJSONFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed to open the ALLOCATION balance snapshot: %v", err))
	}
	defer allocationSnapshotJSONFile.Close()

	var allocationBalanceMap map[string]string
	allocationBalanceMapByteValue, err := ioutil.ReadAll(allocationSnapshotJSONFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read the ALLOCATION balance snapshot: %v", err))
	}

	json.Unmarshal(allocationBalanceMapByteValue, &allocationBalanceMap)
	for key, val := range allocationBalanceMap {
		if !common.IsHexAddress(key) {
			panic(fmt.Sprintf("Invalid address: %v", key))
		}
		address := common.HexToAddress(key)

		dnero, success := new(big.Int).SetString(val, 10)
		if !success {
			panic(fmt.Sprintf("Failed to parse DneroWei amount: %v", val))
		}
		dtoken := new(big.Int).Mul(initDTokenToDneroRatio, dnero)
		acc := &types.Account{
			Address:  address,
			Root:     common.Hash{},
			CodeHash: types.EmptyCodeHash,
			Balance: types.Coins{
				DneroWei: dnero,
				DTokenWei: dtoken,
			},
		}
		sv.SetAccount(acc.Address, acc)
		//logger.Infof("address: %v, dnero: %v, dtoken: %v", strings.ToLower(address.String()), dnero, dtoken)
	}

	return sv
}

func performInitialStakeDeposit(stakeDepositFilePath string, genesisHeight uint64, sv *state.StoreView) *core.ValidatorCandidatePool {
	var stakeDeposits []StakeDeposit
	stakeDepositFile, err := os.Open(stakeDepositFilePath)
	stakeDepositByteValue, err := ioutil.ReadAll(stakeDepositFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read initial stake deposit file: %v", err))
	}

	json.Unmarshal(stakeDepositByteValue, &stakeDeposits)
	vcp := &core.ValidatorCandidatePool{}
	for _, stakeDeposit := range stakeDeposits {
		if !common.IsHexAddress(stakeDeposit.Source) {
			panic(fmt.Sprintf("Invalid source address: %v", stakeDeposit.Source))
		}
		if !common.IsHexAddress(stakeDeposit.Holder) {
			panic(fmt.Sprintf("Invalid holder address: %v", stakeDeposit.Holder))
		}
		sourceAddress := common.HexToAddress(stakeDeposit.Source)
		holderAddress := common.HexToAddress(stakeDeposit.Holder)
		stakeAmount, success := new(big.Int).SetString(stakeDeposit.Amount, 10)
		if !success {
			panic(fmt.Sprintf("Failed to parse Stake amount: %v", stakeDeposit.Amount))
		}

		sourceAccount := sv.GetAccount(sourceAddress)
		if sourceAccount == nil {
			panic(fmt.Sprintf("Failed to retrieve account for source address: %v", sourceAddress))
		}
		if sourceAccount.Balance.DneroWei.Cmp(stakeAmount) < 0 {
			panic(fmt.Sprintf("The source account %v does NOT have sufficient balance for stake deposit. DneroWeiBalance = %v, StakeAmount = %v",
				sourceAddress, sourceAccount.Balance.DneroWei, stakeDeposit.Amount))
		}
		err := vcp.DepositStake(sourceAddress, holderAddress, stakeAmount, genesisHeight)
		if err != nil {
			panic(fmt.Sprintf("Failed to deposit stake, err: %v", err))
		}

		stake := types.Coins{
			DneroWei: stakeAmount,
			DTokenWei: new(big.Int).SetUint64(0),
		}
		sourceAccount.Balance = sourceAccount.Balance.Minus(stake)
		sv.SetAccount(sourceAddress, sourceAccount)
	}

	sv.UpdateValidatorCandidatePool(vcp)

	hl := &types.HeightList{}
	hl.Append(genesisHeight)
	sv.UpdateStakeTransactionHeightList(hl)

	return vcp
}

func proveVCP(sv *state.StoreView) (*core.VCPProof, error) {
	vp := &core.VCPProof{}
	vcpKey := state.ValidatorCandidatePoolKey()
	err := sv.ProveVCP(vcpKey, vp)
	return vp, err
}

// writeGenesisSnapshot writes genesis snapshot to file system.
func writeGenesisSnapshot(sv *state.StoreView, metadata *core.SnapshotMetadata, genesisSnapshotFilePath string) error {
	file, err := os.Create(genesisSnapshotFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = core.WriteMetadata(writer, metadata)
	if err != nil {
		return err
	}
	writeStoreView(sv, true, writer)
	return err
}

func writeStoreView(sv *state.StoreView, needAccountStorage bool, writer *bufio.Writer) {
	height := core.Itobytes(sv.Height())
	err := core.WriteRecord(writer, []byte{core.SVStart}, height)
	if err != nil {
		panic(err)
	}
	sv.GetStore().Traverse(nil, func(k, v common.Bytes) bool {
		err = core.WriteRecord(writer, k, v)
		if err != nil {
			panic(err)
		}
		return true
	})
	err = core.WriteRecord(writer, []byte{core.SVEnd}, height)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

func sanityChecks(sv *state.StoreView) error {
	dneroWeiTotal := new(big.Int).SetUint64(0)
	dtokenWeiTotal := new(big.Int).SetUint64(0)

	vcpAnalyzed := false
	sv.GetStore().Traverse(nil, func(key, val common.Bytes) bool {
		if bytes.Compare(key, state.ValidatorCandidatePoolKey()) == 0 {
			var vcp core.ValidatorCandidatePool
			err := rlp.DecodeBytes(val, &vcp)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode VCP: %v", err))
			}
			for _, sc := range vcp.SortedCandidates {
				logger.Infof("--------------------------------------------------------")
				logger.Infof("Validator Candidate: %v, totalStake  = %v", sc.Holder, sc.TotalStake())
				for _, stake := range sc.Stakes {
					dneroWeiTotal = new(big.Int).Add(dneroWeiTotal, stake.Amount)
					logger.Infof("     Stake: source = %v, stakeAmount = %v", stake.Source, stake.Amount)
				}
				logger.Infof("--------------------------------------------------------")
			}
			vcpAnalyzed = true
		} else if bytes.Compare(key, state.StakeTransactionHeightListKey()) == 0 {
			var hl types.HeightList
			err := rlp.DecodeBytes(val, &hl)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode Height List: %v", err))
			}
			if len(hl.Heights) != 1 {
				panic(fmt.Sprintf("The genesis height list should contain only one height: %v", hl.Heights))
			}
			if hl.Heights[0] != uint64(0) {
				panic(fmt.Sprintf("Only height 0 should be in the genesis height list"))
			}
		} else { // regular account
			var account types.Account
			err := rlp.DecodeBytes(val, &account)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode Account: %v", err))
			}

			dneroWei := account.Balance.DneroWei
			dtokenWei := account.Balance.DTokenWei
			dneroWeiTotal = new(big.Int).Add(dneroWeiTotal, dneroWei)
			dtokenWeiTotal = new(big.Int).Add(dtokenWeiTotal, dtokenWei)

			logger.Infof("Account: %v, DneroWei = %v, DTokenWei = %v", account.Address, dneroWei, dtokenWei)
		}
		return true
	})

	// Check #1: VCP analyzed
	vcpProof, err := proveVCP(sv)
	if err != nil {
		panic(fmt.Sprintf("Failed to get VCP proof from storeview"))
	}
	_, _, err = trie.VerifyProof(sv.Hash(), state.ValidatorCandidatePoolKey(), vcpProof)
	if err != nil {
		panic(fmt.Sprintf("Failed to verify VCP proof in storeview"))
	}
	if !vcpAnalyzed {
		return fmt.Errorf("VCP not detected in the genesis file")
	}

	// Check #2: Sum(DneroWei) + Sum(Stake) == 1 * 10^9 * 10^18
	oneBillion := new(big.Int).SetUint64(1000000000)
	fiveBillion := new(big.Int).Mul(new(big.Int).SetUint64(5), oneBillion)
	ten18 := new(big.Int).SetUint64(1000000000000000000)

	expectedDneroWeiTotal := new(big.Int).Mul(oneBillion, ten18)
	if expectedDneroWeiTotal.Cmp(dneroWeiTotal) != 0 {
		return fmt.Errorf("Unmatched DneroWei total: expected = %v, calculated = %v", expectedDneroWeiTotal, dneroWeiTotal)
	}
	logger.Infof("Expected   DneroWei total = %v", expectedDneroWeiTotal)
	logger.Infof("Calculated DneroWei total = %v", dneroWeiTotal)

	// Check #3: Sum(DTokenWei) == 5 * 10^9 * 10^18
	expectedDTokenWeiTotal := new(big.Int).Mul(fiveBillion, ten18)
	if expectedDTokenWeiTotal.Cmp(dtokenWeiTotal) != 0 {
		return fmt.Errorf("Unmatched DTokenWei total: expected = %v, calculated = %v", expectedDTokenWeiTotal, dtokenWeiTotal)
	}
	logger.Infof("Expected   DTokenWei total = %v", expectedDTokenWeiTotal)
	logger.Infof("Calculated DTokenWei total = %v", dtokenWeiTotal)

	return nil
}

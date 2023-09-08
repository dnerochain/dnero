package core

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/common/result"
	"github.com/dnerochain/dnero/crypto"
	"github.com/dnerochain/dnero/crypto/bls"
	"github.com/dnerochain/dnero/rlp"
)

//
// ------- AggregatedVotes ------- //
//

// AggregatedVotes represents votes on a block.
type AggregatedVotes struct {
	Block      common.Hash    // Hash of the block.
	Scp        common.Hash    // Hash of sentry candidate pool.
	Multiplies []uint32       // Multiplies of each signer.
	Signature  *bls.Signature // Aggregated signiature.
}

func NewAggregateVotes(block common.Hash, scp *SentryCandidatePool) *AggregatedVotes {
	return &AggregatedVotes{
		Block:      block,
		Scp:        scp.Hash(),
		Multiplies: make([]uint32, scp.WithStake().Len()),
		Signature:  bls.NewAggregateSignature(),
	}
}

func (a *AggregatedVotes) String() string {
	return fmt.Sprintf("AggregatedVotes{Block: %s, Scp: %s,  Multiplies: %v}", a.Block.Hex(), a.Scp.Hex(), a.Multiplies)
}

// signBytes returns the bytes to be signed.
func (a *AggregatedVotes) signBytes() common.Bytes {
	tmp := &AggregatedVotes{
		Block: a.Block,
		Scp:   a.Scp,
	}
	b, _ := rlp.EncodeToBytes(tmp)
	return b
}

// Sign adds signer's signature. Returns false if signer has already signed.
func (a *AggregatedVotes) Sign(key *bls.SecretKey, signerIdx int) bool {
	if a.Multiplies[signerIdx] > 0 {
		// Already signed, do nothing.
		return false
	}

	a.Multiplies[signerIdx] = 1
	a.Signature.Aggregate(key.Sign(a.signBytes()))
	return true
}

// Merge creates a new aggregation that combines two vote sets. Returns nil, nil if input vote
// is a subset of current vote.
func (a *AggregatedVotes) Merge(b *AggregatedVotes) (*AggregatedVotes, error) {
	if a.Block != b.Block || a.Scp != b.Scp {
		return nil, errors.New("Cannot merge incompatible votes")
	}
	newMultiplies := make([]uint32, len(a.Multiplies))
	isSubset := true
	for i := 0; i < len(a.Multiplies); i++ {
		newMultiplies[i] = a.Multiplies[i] + b.Multiplies[i]
		if newMultiplies[i] < a.Multiplies[i] || newMultiplies[i] < b.Multiplies[i] {
			return nil, errors.New("Signiature multipliers overflowed")
		}
		if a.Multiplies[i] == 0 && b.Multiplies[i] != 0 {
			isSubset = false
		}
	}
	if isSubset {
		// The other vote is a subset of current vote
		return nil, nil
	}
	newSig := a.Signature.Copy()
	newSig.Aggregate(b.Signature)
	return &AggregatedVotes{
		Block:      a.Block,
		Scp:        a.Scp,
		Multiplies: newMultiplies,
		Signature:  newSig,
	}, nil
}

// Abs returns the number of voted sentrys in the vote
func (a *AggregatedVotes) Abs() int {
	ret := 0
	for i := 0; i < len(a.Multiplies); i++ {
		if a.Multiplies[i] != 0 {
			ret += 1
		}
	}
	return ret
}

// Pick selects better vote from two votes.
func (a *AggregatedVotes) Pick(b *AggregatedVotes) (*AggregatedVotes, error) {
	if a.Block != b.Block || a.Scp != b.Scp {
		return nil, errors.New("Cannot compare incompatible votes")
	}
	if b.Abs() > a.Abs() {
		return b, nil
	}
	return a, nil
}

// Validate verifies the voteset.
func (a *AggregatedVotes) Validate(scp *SentryCandidatePool) result.Result {
	if scp.Hash() != a.Scp {
		return result.Error("scp hash mismatch: scp.Hash(): %s, vote.Scp: %s", scp.Hash().Hex(), a.Scp.Hex())
	}
	if len(a.Multiplies) != scp.WithStake().Len() {
		return result.Error("multiplies size %d is not equal to scp size %d", len(a.Multiplies), scp.WithStake().Len())
	}
	if a.Signature == nil {
		return result.Error("signature cannot be nil")
	}
	pubKeys := scp.WithStake().PubKeys()
	aggPubkey := bls.AggregatePublicKeysVec(pubKeys, a.Multiplies)
	if !a.Signature.Verify(a.signBytes(), aggPubkey) {
		return result.Error("signature verification failed")
	}
	return result.OK
}

// Copy clones the aggregated votes
func (a *AggregatedVotes) Copy() *AggregatedVotes {
	clone := &AggregatedVotes{
		Block: a.Block,
		Scp:   a.Scp,
	}
	if a.Multiplies != nil {
		clone.Multiplies = make([]uint32, len(a.Multiplies))
		copy(clone.Multiplies, a.Multiplies)
	}
	if a.Signature != nil {
		clone.Signature = a.Signature.Copy()
	}

	return clone
}

//
// ------- SentryCandidatePool ------- //
//

var (
	MinSentryStakeDeposit *big.Int

	//MinSentryStakeDeposit1000 *big.Int
)

func init() {
	// Each stake deposit needs to be at least 2,000 Dnero
	MinSentryStakeDeposit = new(big.Int).Mul(new(big.Int).SetUint64(2000), new(big.Int).SetUint64(1e18))

	// Lowering the sentry stake threshold to 1,000 Dnero
	//MinSentryStakeDeposit1000 = new(big.Int).Mul(new(big.Int).SetUint64(1000), new(big.Int).SetUint64(1e18))

}

type SentryCandidatePool struct {
	SortedSentrys []*Sentry // Sentrys sorted by holder address.
}

// NewSentryCandidatePool creates a new instance of SentryCandidatePool.
func NewSentryCandidatePool() *SentryCandidatePool {
	return &SentryCandidatePool{
		SortedSentrys: []*Sentry{},
	}
}

// Add inserts sentry into the pool; returns false if sentry is already added.
func (scp *SentryCandidatePool) Add(g *Sentry) bool {
	k := sort.Search(scp.Len(), func(i int) bool {
		return bytes.Compare(scp.SortedSentrys[i].Holder.Bytes(), g.Holder.Bytes()) >= 0
	})

	if k == scp.Len() {
		scp.SortedSentrys = append(scp.SortedSentrys, g)
		return true
	}

	// Sentry is already added.
	if scp.SortedSentrys[k].Holder == g.Holder {
		return false
	}
	scp.SortedSentrys = append(scp.SortedSentrys, nil)
	copy(scp.SortedSentrys[k+1:], scp.SortedSentrys[k:])
	scp.SortedSentrys[k] = g
	return true
}

// Remove removes a sentry from the pool; returns false if sentry is not found.
func (scp *SentryCandidatePool) Remove(g common.Address) bool {
	k := sort.Search(scp.Len(), func(i int) bool {
		return bytes.Compare(scp.SortedSentrys[i].Holder.Bytes(), g.Bytes()) >= 0
	})

	if k == scp.Len() || bytes.Compare(scp.SortedSentrys[k].Holder.Bytes(), g.Bytes()) != 0 {
		return false
	}
	scp.SortedSentrys = append(scp.SortedSentrys[:k], scp.SortedSentrys[k+1:]...)
	return true
}

// Contains checks if given address is in the pool.
func (scp *SentryCandidatePool) Contains(g common.Address) bool {
	k := sort.Search(scp.Len(), func(i int) bool {
		return bytes.Compare(scp.SortedSentrys[i].Holder.Bytes(), g.Bytes()) >= 0
	})

	if k == scp.Len() || scp.SortedSentrys[k].Holder != g {
		return false
	}
	return true
}

// WithStake returns a new pool with withdrawn sentrys filtered out.
func (scp *SentryCandidatePool) WithStake() *SentryCandidatePool {
	ret := NewSentryCandidatePool()
	for _, g := range scp.SortedSentrys {
		// Skip if sentry dons't have non-withdrawn stake
		hasStake := false
		for _, stake := range g.Stakes {
			if !stake.Withdrawn {
				hasStake = true
				break
			}
		}
		if !hasStake {
			continue
		}

		ret.Add(g)
	}
	return ret
}

// GetWithHolderAddress returns the sentry node correspond to the stake holder in the pool. Returns nil if not found.
func (scp *SentryCandidatePool) GetWithHolderAddress(addr common.Address) *Sentry {
	for _, g := range scp.SortedSentrys {
		if g.Holder == addr {
			return g
		}
	}
	return nil
}

// Index returns index of a public key in the pool. Returns -1 if not found.
func (scp *SentryCandidatePool) Index(pubkey *bls.PublicKey) int {
	for i, g := range scp.SortedSentrys {
		if pubkey.Equals(g.Pubkey) {
			return i
		}
	}
	return -1
}

// PubKeys exports sentrys' public keys.
func (scp *SentryCandidatePool) PubKeys() []*bls.PublicKey {
	ret := make([]*bls.PublicKey, scp.Len())
	for i, g := range scp.SortedSentrys {
		ret[i] = g.Pubkey
	}
	return ret
}

// Implements sort.Interface for Sentrys based on
// the Address field.
func (scp *SentryCandidatePool) Len() int {
	return len(scp.SortedSentrys)
}
func (scp *SentryCandidatePool) Swap(i, j int) {
	scp.SortedSentrys[i], scp.SortedSentrys[j] = scp.SortedSentrys[j], scp.SortedSentrys[i]
}
func (scp *SentryCandidatePool) Less(i, j int) bool {
	return bytes.Compare(scp.SortedSentrys[i].Holder.Bytes(), scp.SortedSentrys[j].Holder.Bytes()) < 0
}

// Hash calculates the hash of scp.
func (scp *SentryCandidatePool) Hash() common.Hash {
	raw, err := rlp.EncodeToBytes(scp)
	if err != nil {
		logger.Panic(err)
	}
	return crypto.Keccak256Hash(raw)
}

func (scp *SentryCandidatePool) DepositStake(source common.Address, holder common.Address, amount *big.Int, pubkey *bls.PublicKey, blockHeight uint64) (err error) {
	minSentryStake := MinSentryStakeDeposit
	//if blockHeight >= common.HeightLowerGNStakeThresholdTo1000 { //StakeDeposit Fork Removed
		//minSentryStake = MinSentryStakeDeposit1000
	//}
	if amount.Cmp(minSentryStake) < 0 {
		return fmt.Errorf("Insufficient stake: %v", amount)
	}

	matchedHolderFound := false
	for _, candidate := range scp.SortedSentrys {
		if candidate.Holder == holder {
			matchedHolderFound = true
			err = candidate.depositStake(source, amount)
			if err != nil {
				return err
			}
			break
		}
	}

	if !matchedHolderFound {
		newSentry := &Sentry{
			StakeHolder: NewStakeHolder(holder, []*Stake{NewStake(source, amount)}),
			Pubkey:      pubkey,
		}
		scp.Add(newSentry)
	}
	return nil
}

func (scp *SentryCandidatePool) WithdrawStake(source common.Address, holder common.Address, currentHeight uint64) error {
	matchedHolderFound := false
	for _, g := range scp.SortedSentrys {
		if g.Holder == holder {
			matchedHolderFound = true
			_, err := g.withdrawStake(source, currentHeight)
			if err != nil {
				return err
			}
			break
		}
	}

	if !matchedHolderFound {
		return fmt.Errorf("No matched stake holder address found: %v", holder)
	}
	return nil
}

func (scp *SentryCandidatePool) ReturnStakes(currentHeight uint64) []*Stake {
	returnedStakes := []*Stake{}

	// need to iterate in the reverse order, since we may delete elemements
	// from the slice while iterating through it
	for cidx := scp.Len() - 1; cidx >= 0; cidx-- {
		g := scp.SortedSentrys[cidx]
		numStakeSources := len(g.Stakes)
		for sidx := numStakeSources - 1; sidx >= 0; sidx-- { // similar to the outer loop, need to iterate in the reversed order
			stake := g.Stakes[sidx]
			if (stake.Withdrawn) && (currentHeight >= stake.ReturnHeight) {
				logger.Printf("Stake to be returned: source = %v, amount = %v", stake.Source, stake.Amount)
				source := stake.Source
				returnedStake, err := g.returnStake(source, currentHeight)
				if err != nil {
					logger.Errorf("Failed to return stake: %v, error: %v", source, err)
					continue
				}
				returnedStakes = append(returnedStakes, returnedStake)
			}
		}

		if len(g.Stakes) == 0 { // the candidate's stake becomes zero, no need to keep track of the candidate anymore
			scp.Remove(g.Holder)
		}
	}
	return returnedStakes
}

//
// ------- Sentry ------- //
//

type Sentry struct {
	*StakeHolder
	Pubkey *bls.PublicKey `json:"-"`
}

func (g *Sentry) String() string {
	return fmt.Sprintf("{holder: %v, pubkey: %v, stakes :%v}", g.Holder, g.Pubkey.String(), g.Stakes)
}

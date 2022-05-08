package core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/crypto/bls"
	"github.com/dnerochain/dnero/rlp"

	"github.com/dnerochain/dnero/crypto"
)

func createTestSentryPool(size int) (*SentryCandidatePool, map[common.Address]*bls.SecretKey) {
	pool := NewSentryCandidatePool()
	sks := make(map[common.Address]*bls.SecretKey)
	for i := 0; i < size; i++ {
		_, pub, _ := crypto.GenerateKeyPair()
		blsKey, _ := bls.RandKey()
		g := &Sentry{
			StakeHolder: &StakeHolder{
				Holder: pub.Address(),
				Stakes: []*Stake{&Stake{
					Source:       pub.Address(),
					Amount:       MinSentryStakeDeposit,
					Withdrawn:    false,
					ReturnHeight: 99999999999,
				}},
			},
			Pubkey: blsKey.PublicKey(),
		}
		pool.Add(g)
		sks[g.Holder] = blsKey
	}
	return pool, sks
}

func isSorted(pl *SentryCandidatePool) bool {
	g := pl.SortedSentrys[0]
	for i := 1; i < pl.Len(); i++ {
		if bytes.Compare(g.Holder.Bytes(), pl.SortedSentrys[i].Holder.Bytes()) >= 0 {
			return false
		}
	}
	return true
}

func TestSentryPool(t *testing.T) {
	require := require.New(t)

	pool, _ := createTestSentryPool(10)

	// Should be sorted.
	if !isSorted(pool) {
		t.Fatal("Sentry pool is not sorted")
	}

	// Should not add duplicate.
	newSentry := &Sentry{
		StakeHolder: &StakeHolder{
			Holder: pool.SortedSentrys[3].Holder,
		},
	}
	if pool.Add(newSentry) {
		t.Fatal("Should not add duplicate sentry")
	}

	// Should add new sentry.
	_, pub, _ := crypto.GenerateKeyPair()
	blsKey, _ := bls.RandKey()
	g := &Sentry{
		StakeHolder: &StakeHolder{
			Holder: pub.Address(),
			Stakes: []*Stake{&Stake{
				Source:       pub.Address(),
				Amount:       MinSentryStakeDeposit,
				Withdrawn:    false,
				ReturnHeight: 99999999999,
			}},
		},
		Pubkey: blsKey.PublicKey(),
	}
	if !pool.Add(g) || pool.Len() != 11 {
		t.Fatal("Should add new sentry")
	}
	if !isSorted(pool) {
		t.Fatal("Should be sorted after add")
	}

	// Should remove sentry.
	toRemove := pool.SortedSentrys[5].Holder
	toRemoveBlsPub := pool.SortedSentrys[5].Pubkey
	if !pool.Remove(toRemove) || pool.Len() != 10 {
		t.Fatal("Should remove sentry")
	}
	if !isSorted(pool) {
		t.Fatal("Should be sorted after remove")
	}

	// Should return false when removing non-existent sentry.
	if pool.Remove(toRemove) || pool.Len() != 10 {
		t.Fatal("Should not remove non-existent sentry")
	}

	// Should return -1 for removed sentry.
	require.Equal(-1, pool.Index(toRemoveBlsPub), "Should return -1 for removed sentry")

	toWithdrawnPub := pool.SortedSentrys[3].Pubkey
	nextPub := pool.SortedSentrys[4].Pubkey
	require.Equal(3, pool.WithStake().Index(toWithdrawnPub))
	require.Equal(4, pool.WithStake().Index(nextPub))
	pool.SortedSentrys[3].Stakes[0].Withdrawn = true
	// Should return -1 for withdrawn sentry.
	require.Equal(-1, pool.WithStake().Index(toWithdrawnPub))
	// Should skip withdrawn sentry.
	require.Equal(3, pool.WithStake().Index(nextPub))
}

func TestAggregateVote(t *testing.T) {
	pool, sks := createTestSentryPool(10)

	bh := common.BytesToHash([]byte{12})
	vote1 := NewAggregateVotes(bh, pool)

	g1 := pool.SortedSentrys[0].Holder

	// Sentry 1 signs a vote.
	success := vote1.Sign(sks[g1], 0)
	if !success {
		t.Fatal("Should sign")
	}
	if res := vote1.Validate(pool); res.IsError() {
		t.Fatal("Should validate", res.Message)
	}

	// Sentry 2 signs a vote.
	vote2 := NewAggregateVotes(bh, pool)
	g2 := pool.SortedSentrys[1].Holder
	success = vote2.Sign(sks[g2], 1)
	if !success {
		t.Fatal("Should sign")
	}
	if res := vote2.Validate(pool); res.IsError() {
		t.Fatal("Should validate", res.Message)
	}

	// Should merge two votes.
	vote12, err := vote1.Merge(vote2)
	if err != nil {
		t.Fatalf("Failed to merge votes: %s", err.Error())
	}
	if res := vote12.Validate(pool); res.IsError() {
		t.Fatal("Should validate", res.Message)
	}

	// Should not merge votes that is a subset of current vote.
	res, err := vote12.Merge(vote2)
	if err != nil || res != nil {
		t.Fatalf("Should not merge votes that is subset")
	}
	res, err = vote12.Merge(NewAggregateVotes(bh, pool))
	if err != nil || res != nil {
		t.Fatalf("Should not merge votes that is subset")
	}
	res, err = vote12.Merge(vote12)
	if err != nil || res != nil {
		t.Fatalf("Should not merge votes that is subset")
	}
}

func TestAggregateVoteEncoding(t *testing.T) {
	require := require.New(t)

	pool, sks := createTestSentryPool(10)

	bh := common.BytesToHash([]byte{12})
	vote1 := NewAggregateVotes(bh, pool)

	g1 := pool.SortedSentrys[0].Holder

	// Sentry 1 signs a vote.
	success := vote1.Sign(sks[g1], 0)
	require.True(success, "Should sign")

	raw, err := rlp.EncodeToBytes(vote1)
	require.Nil(err)

	vote2 := &AggregatedVotes{}
	err = rlp.DecodeBytes(raw, vote2)
	require.Nil(err)
}

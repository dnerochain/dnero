package blockchain

import (
	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/core"
	"github.com/dnerochain/dnero/crypto"
	"github.com/dnerochain/dnero/ledger/types"
	"github.com/dnerochain/dnero/store"
)

// txIndexKey constructs the DB key for the given transaction hash.
func txIndexKey(hash common.Hash) common.Bytes {
	return append(common.Bytes("tx/"), hash[:]...)
}

// TxIndexEntry is a positional metadata to help looking up a transaction given only its hash.
type TxIndexEntry struct {
	BlockHash   common.Hash
	BlockHeight uint64
	Index       uint64
}

// AddTxsToIndex adds transactions in given block to index.
func (ch *Chain) AddTxsToIndex(block *core.ExtendedBlock, force bool) {
	for idx, tx := range block.Txs {
		txIndexEntry := TxIndexEntry{
			BlockHash:   block.Hash(),
			BlockHeight: block.Height,
			Index:       uint64(idx),
		}
		txHash := crypto.Keccak256Hash(tx)
		key := txIndexKey(txHash)

		if !force {
			// Check if TX with given hash exists in DB.
			err := ch.store.Get(key, &TxIndexEntry{})
			if err != store.ErrKeyNotFound {
				continue
			}
		}

		err := ch.store.Put(key, txIndexEntry)
		if err != nil {
			logger.Panic(err)
		}
	}
}

// FindTxByHash looks up transaction by hash and additionally returns the containing block.
func (ch *Chain) FindTxByHash(hash common.Hash) (tx common.Bytes, block *core.ExtendedBlock, founded bool) {
	txIndexEntry := &TxIndexEntry{}
	err := ch.store.Get(txIndexKey(hash), txIndexEntry)
	if err != nil {
		if err != store.ErrKeyNotFound {
			logger.Error(err)
		}
		return nil, nil, false
	}
	block, err = ch.FindBlock(txIndexEntry.BlockHash)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, nil, false
		}
		logger.Panic(err)
	}
	return block.Txs[txIndexEntry.Index], block, true
}

// ---------------- Tx Receipts ---------------

// txReceiptKey constructs the DB key for the given transaction hash.
func txReceiptKey(hash common.Hash) common.Bytes {
	return append(common.Bytes("txr/"), hash[:]...)
}

// TxReceiptEntry records smart contract Tx execution result.
type TxReceiptEntry struct {
	TxHash          common.Hash
	Logs            []*types.Log
	EvmRet          common.Bytes
	ContractAddress common.Address
	GasUsed         uint64
	EvmErr          string
}

// AddTxReceipt adds transaction receipt.
func (ch *Chain) AddTxReceipt(tx types.Tx, logs []*types.Log, evmRet common.Bytes,
	contractAddr common.Address, gasUsed uint64, evmErr error) {
	raw, err := types.TxToBytes(tx)
	if err != nil {
		// Should never happen
		logger.Panic(err)
	}
	txHash := crypto.Keccak256Hash(raw)
	errStr := ""
	if evmErr != nil {
		errStr = evmErr.Error()
	}
	txReceiptEntry := TxReceiptEntry{
		TxHash:          txHash,
		Logs:            logs,
		EvmRet:          evmRet,
		ContractAddress: contractAddr,
		GasUsed:         gasUsed,
		EvmErr:          errStr,
	}
	key := txReceiptKey(txHash)

	err = ch.store.Put(key, txReceiptEntry)
	if err != nil {
		logger.Panic(err)
	}
}

// FindTxReceiptByHash looks up transaction receipt by hash.
func (ch *Chain) FindTxReceiptByHash(hash common.Hash) (*TxReceiptEntry, bool) {
	txReceiptEntry := &TxReceiptEntry{}

	key := txReceiptKey(hash)

	err := ch.store.Get(key, txReceiptEntry)

	if err != nil {
		if err != store.ErrKeyNotFound {
			logger.Error(err)
		}
		return nil, false
	}
	return txReceiptEntry, true
}

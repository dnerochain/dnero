package blockchain

import (
	"fmt"
	"math/big"

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

		ch.insertEthTxHash(block, tx, &txIndexEntry)
	}
}

// Index the ETH smart contract transactions, using the ETH tx hash as the key
func (ch *Chain) insertEthTxHash(block *core.ExtendedBlock, rawTxBytes []byte, txIndexEntry *TxIndexEntry) error {
	ethTxHash, err := CalcEthTxHash(block, rawTxBytes)
	if err != nil {
		return err // skip insertion
	}

	key := txIndexKey(ethTxHash)
	err = ch.store.Put(key, *txIndexEntry)
	if err != nil {
		logger.Panic(err)
	}

	return nil
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

// ---------------- Utils ---------------

func CalcEthTxHash(block *core.ExtendedBlock, rawTxBytes []byte) (common.Hash, error) {
	tx, err := types.TxFromBytes(rawTxBytes)
	if err != nil {
		return common.Hash{}, err
	}

	sctx, ok := tx.(*types.SmartContractTx)
	if !ok {
		return common.Hash{}, fmt.Errorf("not a smart contract transaction") // not a smart contract tx, skip ETH tx insertion
	}

	ethSigningHash := sctx.EthSigningHash(block.ChainID, block.Height)
	err = crypto.ValidateEthSignature(sctx.From.Address, ethSigningHash, sctx.From.Signature)
	if err != nil {
		return common.Hash{}, fmt.Errorf("not an ETH smart contract transaction") // it is a Dnero native smart contract transaction, no need to index it as an EthTxHash
	}

	var toAddress *common.Address
	if (sctx.To.Address != common.Address{}) {
		toAddress = &sctx.To.Address
	}

	r, s, v := crypto.DecodeSignature(sctx.From.Signature)
	chainID := types.MapChainID(block.ChainID, block.Height)
	vPrime := big.NewInt(1).Mul(chainID, big.NewInt(2))
	vPrime = big.NewInt(0).Add(vPrime, big.NewInt(8))
	vPrime = big.NewInt(0).Add(vPrime, v)

	ethTx := types.EthTransaction{
		Nonce:    sctx.From.Sequence - 1, // off-by-one, ETH tx nonce starts from 0, while Dnero tx sequence starts from 1
		GasPrice: sctx.GasPrice,
		Gas:      sctx.GasLimit,
		To:       toAddress,
		Value:    sctx.From.Coins.NoNil().DTokenWei,
		Data:     sctx.Data,
		V:        vPrime,
		R:        r,
		S:        s,
	}

	ethTxHash := ethTx.Hash()

	//ethTxBytes, _ := rlp.EncodeToBytes(ethTx)
	//logger.Debugf("ethTxBytes: %v", hex.EncodeToString(ethTxBytes))
	logger.Debugf("ethTxHash: %v", ethTxHash.Hex())
	logger.Debugf("ethTxHash, nonce: %v, r: %x, s: %x, v: %v", sctx.From.Sequence-1, r, s, vPrime)

	return ethTxHash, nil
}

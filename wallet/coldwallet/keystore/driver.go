package keystore

import (
	"io"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/crypto"
	"github.com/dnerochain/dnero/wallet/types"
)

// Driver abstracts the functionality of the hardware wallet
type Driver interface {
	Status() (string, error)
	Open(device io.ReadWriter, password string) error
	Close() error
	Heartbeat() error
	Derive(path types.DerivationPath) (common.Address, error)
	SignTx(path types.DerivationPath, txrlp common.Bytes) (common.Address, *crypto.Signature, error)
}

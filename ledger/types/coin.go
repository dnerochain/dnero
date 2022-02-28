package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/dnerochain/dnero/common"
)

var (
	Zero    *big.Int
	Hundred *big.Int
)

func init() {
	Zero = big.NewInt(0)
	Hundred = big.NewInt(100)
}

type Coins struct {
	DneroWei *big.Int
	DTokenWei *big.Int
}

type CoinsJSON struct {
	DneroWei *common.JSONBig `json:"dnerowei"`
	DTokenWei *common.JSONBig `json:"dtokenwei"`
}

func NewCoinsJSON(coin Coins) CoinsJSON {
	return CoinsJSON{
		DneroWei: (*common.JSONBig)(coin.DneroWei),
		DTokenWei: (*common.JSONBig)(coin.DTokenWei),
	}
}

func (c CoinsJSON) Coins() Coins {
	return Coins{
		DneroWei: (*big.Int)(c.DneroWei),
		DTokenWei: (*big.Int)(c.DTokenWei),
	}
}

func (c Coins) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewCoinsJSON(c))
}

func (c *Coins) UnmarshalJSON(data []byte) error {
	var a CoinsJSON
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*c = a.Coins()
	return nil
}

// NewCoins is a convenient method for creating small amount of coins.
func NewCoins(dnero int64, dtoken int64) Coins {
	return Coins{
		DneroWei: big.NewInt(dnero),
		DTokenWei: big.NewInt(dtoken),
	}
}

func (coins Coins) String() string {
	return fmt.Sprintf("%v %v, %v %v", coins.DneroWei, DenomDneroWei, coins.DTokenWei, DenomDTokenWei)
}

func (coins Coins) IsValid() bool {
	return coins.IsNonnegative()
}

func (coins Coins) NoNil() Coins {
	dnero := coins.DneroWei
	if dnero == nil {
		dnero = big.NewInt(0)
	}
	dtoken := coins.DTokenWei
	if dtoken == nil {
		dtoken = big.NewInt(0)
	}

	return Coins{
		DneroWei: dnero,
		DTokenWei: dtoken,
	}
}

// CalculatePercentage function calculates amount of coins for the given the percentage
func (coins Coins) CalculatePercentage(percentage uint) Coins {
	c := coins.NoNil()

	p := big.NewInt(int64(percentage))

	dnero := new(big.Int)
	dnero.Mul(c.DneroWei, p)
	dnero.Div(dnero, Hundred)

	dtoken := new(big.Int)
	dtoken.Mul(c.DTokenWei, p)
	dtoken.Div(dtoken, Hundred)

	return Coins{
		DneroWei: dnero,
		DTokenWei: dtoken,
	}
}

// Currently appends an empty coin ...
func (coinsA Coins) Plus(coinsB Coins) Coins {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()

	dnero := new(big.Int)
	dnero.Add(cA.DneroWei, cB.DneroWei)

	dtoken := new(big.Int)
	dtoken.Add(cA.DTokenWei, cB.DTokenWei)

	return Coins{
		DneroWei: dnero,
		DTokenWei: dtoken,
	}
}

func (coins Coins) Negative() Coins {
	c := coins.NoNil()

	dnero := new(big.Int)
	dnero.Neg(c.DneroWei)

	dtoken := new(big.Int)
	dtoken.Neg(c.DTokenWei)

	return Coins{
		DneroWei: dnero,
		DTokenWei: dtoken,
	}
}

func (coinsA Coins) Minus(coinsB Coins) Coins {
	return coinsA.Plus(coinsB.Negative())
}

func (coinsA Coins) IsGTE(coinsB Coins) bool {
	diff := coinsA.Minus(coinsB)
	return diff.IsNonnegative()
}

func (coins Coins) IsZero() bool {
	c := coins.NoNil()
	return c.DneroWei.Cmp(Zero) == 0 && c.DTokenWei.Cmp(Zero) == 0
}

func (coinsA Coins) IsEqual(coinsB Coins) bool {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()
	return cA.DneroWei.Cmp(cB.DneroWei) == 0 && cA.DTokenWei.Cmp(cB.DTokenWei) == 0
}

func (coins Coins) IsPositive() bool {
	c := coins.NoNil()
	return (c.DneroWei.Cmp(Zero) > 0 && c.DTokenWei.Cmp(Zero) >= 0) ||
		(c.DneroWei.Cmp(Zero) >= 0 && c.DTokenWei.Cmp(Zero) > 0)
}

func (coins Coins) IsNonnegative() bool {
	c := coins.NoNil()
	return c.DneroWei.Cmp(Zero) >= 0 && c.DTokenWei.Cmp(Zero) >= 0
}

// ParseCoinAmount parses a string representation of coin amount.
func ParseCoinAmount(in string) (*big.Int, bool) {
	inWei := false
	if len(in) > 3 && strings.EqualFold("wei", in[len(in)-3:]) {
		inWei = true
		in = in[:len(in)-3]
	}

	f, ok := new(big.Float).SetPrec(1024).SetString(in)
	if !ok || f.Sign() < 0 {
		return nil, false
	}

	if !inWei {
		f = f.Mul(f, new(big.Float).SetPrec(1024).SetUint64(1e18))
	}

	ret, _ := f.Int(nil)

	return ret, true
}

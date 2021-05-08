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
	DFuelWei *big.Int
}

type CoinsJSON struct {
	DneroWei *common.JSONBig `json:"dnerowei"`
	DFuelWei *common.JSONBig `json:"dfuelwei"`
}

func NewCoinsJSON(coin Coins) CoinsJSON {
	return CoinsJSON{
		DneroWei: (*common.JSONBig)(coin.DneroWei),
		DFuelWei: (*common.JSONBig)(coin.DFuelWei),
	}
}

func (c CoinsJSON) Coins() Coins {
	return Coins{
		DneroWei: (*big.Int)(c.DneroWei),
		DFuelWei: (*big.Int)(c.DFuelWei),
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
func NewCoins(dnero int64, dfuel int64) Coins {
	return Coins{
		DneroWei: big.NewInt(dnero),
		DFuelWei: big.NewInt(dfuel),
	}
}

func (coins Coins) String() string {
	return fmt.Sprintf("%v %v, %v %v", coins.DneroWei, DenomDneroWei, coins.DFuelWei, DenomDFuelWei)
}

func (coins Coins) IsValid() bool {
	return coins.IsNonnegative()
}

func (coins Coins) NoNil() Coins {
	dnero := coins.DneroWei
	if dnero == nil {
		dnero = big.NewInt(0)
	}
	dfuel := coins.DFuelWei
	if dfuel == nil {
		dfuel = big.NewInt(0)
	}

	return Coins{
		DneroWei: dnero,
		DFuelWei: dfuel,
	}
}

// CalculatePercentage function calculates amount of coins for the given the percentage
func (coins Coins) CalculatePercentage(percentage uint) Coins {
	c := coins.NoNil()

	p := big.NewInt(int64(percentage))

	dnero := new(big.Int)
	dnero.Mul(c.DneroWei, p)
	dnero.Div(dnero, Hundred)

	dfuel := new(big.Int)
	dfuel.Mul(c.DFuelWei, p)
	dfuel.Div(dfuel, Hundred)

	return Coins{
		DneroWei: dnero,
		DFuelWei: dfuel,
	}
}

// Currently appends an empty coin ...
func (coinsA Coins) Plus(coinsB Coins) Coins {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()

	dnero := new(big.Int)
	dnero.Add(cA.DneroWei, cB.DneroWei)

	dfuel := new(big.Int)
	dfuel.Add(cA.DFuelWei, cB.DFuelWei)

	return Coins{
		DneroWei: dnero,
		DFuelWei: dfuel,
	}
}

func (coins Coins) Negative() Coins {
	c := coins.NoNil()

	dnero := new(big.Int)
	dnero.Neg(c.DneroWei)

	dfuel := new(big.Int)
	dfuel.Neg(c.DFuelWei)

	return Coins{
		DneroWei: dnero,
		DFuelWei: dfuel,
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
	return c.DneroWei.Cmp(Zero) == 0 && c.DFuelWei.Cmp(Zero) == 0
}

func (coinsA Coins) IsEqual(coinsB Coins) bool {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()
	return cA.DneroWei.Cmp(cB.DneroWei) == 0 && cA.DFuelWei.Cmp(cB.DFuelWei) == 0
}

func (coins Coins) IsPositive() bool {
	c := coins.NoNil()
	return (c.DneroWei.Cmp(Zero) > 0 && c.DFuelWei.Cmp(Zero) >= 0) ||
		(c.DneroWei.Cmp(Zero) >= 0 && c.DFuelWei.Cmp(Zero) > 0)
}

func (coins Coins) IsNonnegative() bool {
	c := coins.NoNil()
	return c.DneroWei.Cmp(Zero) >= 0 && c.DFuelWei.Cmp(Zero) >= 0
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

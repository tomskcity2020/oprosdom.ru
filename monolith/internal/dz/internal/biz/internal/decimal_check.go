package biz_internal

import (
	"errors"

	"github.com/shopspring/decimal"
)

func (b *BizStruct) DecimalCheck(amount string) error {

	if amount == "" {
		return errors.New("amount empty")
	}

	d, err := decimal.NewFromString(amount)
	if err != nil {
		return errors.New("incorrect amount")
	}

	if d.IsNegative() {
		return errors.New("not valid amount")
	}

	if d.Exponent() < -2 {
		return errors.New("must have 2 digits after dot")
	}

	return nil
}

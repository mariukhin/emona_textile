package model

import "fmt"

type Currency string

const (
	Euro Currency = "eur"
)

var currencies = []Currency{
	Euro,
}

func ParseCurrency(value string) (Currency, error) {
	for _, cur := range currencies {
		if Currency(value) == cur {
			return cur, nil
		}
	}

	return "", fmt.Errorf("unknown currency value %s", value)
}

func (c Currency) Format(amount string) string {
	switch c {
	case Euro:
		return fmt.Sprintf("%s%s", amount, "€")
	default:
		return fmt.Sprintf("%s%s", amount, c)
	}
}

type Cost int

var Zero = Cost(0)

func (c Cost) Add(other Cost) Cost {
	return c + other
}

func (c Cost) Multi(val int) Cost {
	return c * Cost(val)
}

func (c Cost) Format(currency Currency) string {
	return currency.Format(fmt.Sprintf("%d.%03d", c/1000, c%1000))
}

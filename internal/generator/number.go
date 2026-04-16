package generator

import (
	"fmt"
	"math/rand"
)

type NumberGen struct{}

func (g *NumberGen) Name() string { return "Numbers" }
func (g *NumberGen) Description() string {
	return "Field source for IDs, integers, decimals, booleans, and currency amounts"
}
func (g *NumberGen) Kind() Kind { return KindField }
func (g *NumberGen) Fields() []Field {
	return []Field{
		{Name: "id", Desc: "UUID"},
		{Name: "integer", Desc: "Random int 0-10000"},
		{Name: "bigint", Desc: "Random int64 0-9,999,999,999,999"},
		{Name: "float", Desc: "Random float 0-1000 (2 decimals)"},
		{Name: "decimal", Desc: "Random decimal 0-9999.9999 (4 decimals)"},
		{Name: "percentage", Desc: "0-100%"},
		{Name: "currency_amount", Desc: "Currency amount as numeric 0.00-9999.99"},
		{Name: "currency_code", Desc: "ISO currency code"},
		{Name: "currency", Desc: "Formatted currency string (legacy)"},
		{Name: "boolean", Desc: "true or false"},
	}
}

func (g *NumberGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		integer := rng.Intn(10001)
		bigint := rng.Int63n(10_000_000_000_000)
		f := float64(rng.Intn(100000)) / 100.0
		decimal := float64(rng.Intn(100_000_000)) / 10_000.0
		pct := rng.Intn(101)
		cents := rng.Intn(1000000)
		currencyAmount := float64(cents) / 100.0
		boolVal := rng.Intn(2) == 1

		records[i] = map[string]any{
			"id":              genUUID(rng),
			"integer":         integer,
			"bigint":          bigint,
			"float":           f,
			"decimal":         decimal,
			"percentage":      pct,
			"currency_amount": currencyAmount,
			"currency_code":   "USD",
			"currency":        fmt.Sprintf("$%.2f", currencyAmount),
			"boolean":         boolVal,
		}
	}
	return records
}

package generator

import (
	"fmt"
	"math/rand"
)

type NumberGen struct{}

func (g *NumberGen) Name() string        { return "Numbers" }
func (g *NumberGen) Description() string { return "IDs, integers, floats, percentages, currencies" }
func (g *NumberGen) Fields() []Field {
	return []Field{
		{Name: "id", Desc: "UUID"},
		{Name: "integer", Desc: "Random int 0-10000"},
		{Name: "float", Desc: "Random float 0-1000 (2 decimals)"},
		{Name: "percentage", Desc: "0-100%"},
		{Name: "currency", Desc: "$0.00 - $9999.99"},
		{Name: "boolean", Desc: "true or false"},
	}
}

func (g *NumberGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		integer := rng.Intn(10001)
		f := float64(rng.Intn(100000)) / 100.0
		pct := rng.Intn(101)
		cents := rng.Intn(1000000)
		boolVal := rng.Intn(2) == 1

		records[i] = map[string]any{
			"id":         genUUID(rng),
			"integer":    integer,
			"float":      f,
			"percentage": pct,
			"currency":   fmt.Sprintf("$%.2f", float64(cents)/100.0),
			"boolean":    boolVal,
		}
	}
	return records
}

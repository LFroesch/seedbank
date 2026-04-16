package generator

import (
	"math/rand"
)

type AddressGen struct{}

func (g *AddressGen) Name() string        { return "Addresses" }
func (g *AddressGen) Description() string { return "Field source for coherent US addresses" }
func (g *AddressGen) Kind() Kind          { return KindField }
func (g *AddressGen) Fields() []Field {
	return []Field{
		{Name: "street", Desc: "Street address"},
		{Name: "city", Desc: "City (matches state)"},
		{Name: "state", Desc: "2-letter state code"},
		{Name: "zip", Desc: "5-digit zip (matches city)"},
		{Name: "country", Desc: "Country"},
		{Name: "full_address", Desc: "Complete formatted address"},
	}
}

func (g *AddressGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		state, city, zip, _ := pickLocation(rng)
		street := pickStreet(rng)
		full := street + ", " + city + ", " + state + " " + zip

		records[i] = map[string]any{
			"street":       street,
			"city":         city,
			"state":        state,
			"zip":          zip,
			"country":      "United States",
			"full_address": full,
		}
	}
	return records
}

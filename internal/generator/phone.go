package generator

import (
	"fmt"
	"math/rand"
)

type PhoneGen struct{}

func (g *PhoneGen) Name() string { return "Phone Numbers" }
func (g *PhoneGen) Description() string {
	return "Field source for US phone numbers with real area codes"
}
func (g *PhoneGen) Kind() Kind { return KindField }
func (g *PhoneGen) Fields() []Field {
	return []Field{
		{Name: "phone", Desc: "(555) 123-4567 format"},
		{Name: "phone_raw", Desc: "5551234567 digits only"},
		{Name: "phone_intl", Desc: "+1-555-123-4567"},
		{Name: "area_code", Desc: "3-digit area code"},
		{Name: "state", Desc: "State the area code belongs to"},
	}
}

func (g *PhoneGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		// Pick a state to get a coherent area code + state pair
		st := usStates[rng.Intn(len(usStates))]
		area := st.AreaCodes[rng.Intn(len(st.AreaCodes))]
		mid := 200 + rng.Intn(800)
		last := rng.Intn(10000)

		raw := fmt.Sprintf("%s%03d%04d", area, mid, last)
		formatted := fmt.Sprintf("(%s) %03d-%04d", area, mid, last)
		intl := fmt.Sprintf("+1-%s-%03d-%04d", area, mid, last)

		records[i] = map[string]any{
			"phone":      formatted,
			"phone_raw":  raw,
			"phone_intl": intl,
			"area_code":  area,
			"state":      st.Code,
		}
	}
	return records
}

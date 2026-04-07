package generator

import (
	"math/rand"
	"strings"
)

type NameGen struct{}

func (g *NameGen) Name() string        { return "Names" }
func (g *NameGen) Description() string { return "First, last, and full names with gender-matched prefix" }
func (g *NameGen) Fields() []Field {
	return []Field{
		{Name: "first_name", Desc: "First name"},
		{Name: "last_name", Desc: "Last name"},
		{Name: "full_name", Desc: "First + Last"},
		{Name: "prefix", Desc: "Gender-matched prefix"},
		{Name: "username", Desc: "Lowercase first.last + digits"},
	}
}

func (g *NameGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		first, last, prefix, _ := pickGendered(rng)
		firstL := strings.ToLower(first)
		lastL := strings.ToLower(last)

		var username string
		switch rng.Intn(3) {
		case 0:
			username = firstL + "." + lastL
		case 1:
			username = firstL[:1] + lastL + itoa(rng.Intn(99))
		case 2:
			username = firstL + itoa(rng.Intn(999))
		}

		records[i] = map[string]any{
			"first_name": first,
			"last_name":  last,
			"full_name":  first + " " + last,
			"prefix":     prefix,
			"username":   username,
		}
	}
	return records
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

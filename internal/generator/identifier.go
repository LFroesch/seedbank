package generator

import (
	"math/rand"
	"strings"
)

type IdentifierGen struct{}

func (g *IdentifierGen) Name() string { return "Identifiers" }
func (g *IdentifierGen) Description() string {
	return "Field source for UUID and compact identifier variants"
}
func (g *IdentifierGen) Kind() Kind { return KindField }
func (g *IdentifierGen) Fields() []Field {
	return []Field{
		{Name: "uuid", Desc: "Lowercase UUID v4-style"},
		{Name: "uuid_upper", Desc: "Uppercase UUID"},
		{Name: "uuid_compact", Desc: "UUID without dashes"},
		{Name: "short_id", Desc: "12-char compact identifier"},
	}
}

func (g *IdentifierGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		uuid := genUUID(rng)
		compact := strings.ReplaceAll(uuid, "-", "")
		records[i] = map[string]any{
			"uuid":         uuid,
			"uuid_upper":   strings.ToUpper(uuid),
			"uuid_compact": compact,
			"short_id":     compact[:12],
		}
	}
	return records
}

package generator

import (
	"math/rand"
	"strings"
)

// MixGenerator combines fields from multiple generators into merged records.
type MixGenerator struct {
	Gens    []Generator
	Fields_ []Field // combined field list
}

func (g *MixGenerator) Name() string        { return "Custom Mix" }
func (g *MixGenerator) Description() string { return "Combined fields from multiple generators" }
func (g *MixGenerator) Kind() Kind          { return KindField }
func (g *MixGenerator) Fields() []Field     { return g.Fields_ }

func (g *MixGenerator) Generate(count int, rng *rand.Rand) []map[string]any {
	// Generate from each sub-generator
	subRecords := make([][]map[string]any, len(g.Gens))
	for i, gen := range g.Gens {
		subRecords[i] = gen.Generate(count, rng)
	}

	// Merge into combined records with prefixed keys
	records := make([]map[string]any, count)
	for i := range records {
		rec := make(map[string]any)
		for gi, gen := range g.Gens {
			prefix := genPrefix(gen.Name())
			if i < len(subRecords[gi]) {
				for _, f := range gen.Fields() {
					rec[prefix+"."+f.Name] = subRecords[gi][i][f.Name]
				}
			}
		}
		records[i] = rec
	}
	return records
}

// genPrefix turns "Person (Linked)" into "person"
func genPrefix(name string) string {
	if idx := strings.Index(name, "("); idx > 0 {
		name = strings.TrimSpace(name[:idx])
	}
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// BuildMixFields creates combined field list from selected generators.
func BuildMixFields(gens []Generator) []Field {
	var fields []Field
	for _, gen := range gens {
		prefix := genPrefix(gen.Name())
		for _, f := range gen.Fields() {
			fields = append(fields, Field{
				Name: prefix + "." + f.Name,
				Desc: gen.Name() + " — " + f.Desc,
			})
		}
	}
	return fields
}

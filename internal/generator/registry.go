package generator

import "math/rand"

// Field represents a single field definition within a generator.
type Field struct {
	Name string
	Desc string
}

// Generator produces fake data records.
type Generator interface {
	Name() string
	Description() string
	Fields() []Field
	Generate(count int, rng *rand.Rand) []map[string]any
}

// Registry holds all available generators.
var Registry []Generator

// Register adds a generator to the global registry.
func Register(g Generator) {
	Registry = append(Registry, g)
}

func init() {
	// Linked generators (coherent records) first
	Register(&PersonGen{})
	Register(&ProductGen{})
	// Individual field generators
	Register(&NameGen{})
	Register(&LoremGen{})
	Register(&AddressGen{})
	Register(&EmailGen{})
	Register(&PhoneGen{})
	Register(&PhotoGen{})
	Register(&NumberGen{})
	Register(&CompanyGen{})
	Register(&DateGen{})
	Register(&ColorGen{})
}

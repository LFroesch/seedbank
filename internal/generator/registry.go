package generator

import "math/rand"

type Kind string

const (
	KindRecord Kind = "record"
	KindField  Kind = "field"
)

// Field represents a single field definition within a generator.
type Field struct {
	Name string
	Desc string
}

// Generator produces fake data records.
type Generator interface {
	Name() string
	Description() string
	Kind() Kind
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
	// Coherent record builders first, then field sources.
	Register(&PersonGen{})
	Register(&ProductGen{})
	Register(&NameGen{})
	Register(&LoremGen{})
	Register(&AddressGen{})
	Register(&EmailGen{})
	Register(&PhoneGen{})
	Register(&PhotoGen{})
	Register(&NumberGen{})
	Register(&IdentifierGen{})
	Register(&NetworkGen{})
	Register(&CompanyGen{})
	Register(&DateGen{})
	Register(&ColorGen{})
}

func Lookup(name string) Generator {
	needle := normalizeName(name)
	for _, g := range Registry {
		for _, alias := range aliasesFor(g) {
			if normalizeName(alias) == needle {
				return g
			}
		}
	}
	return nil
}

func Grouped() map[Kind][]Generator {
	grouped := map[Kind][]Generator{
		KindRecord: {},
		KindField:  {},
	}
	for _, g := range Registry {
		grouped[g.Kind()] = append(grouped[g.Kind()], g)
	}
	return grouped
}

func KindLabel(kind Kind) string {
	switch kind {
	case KindRecord:
		return "Coherent record builders"
	case KindField:
		return "Field sources"
	default:
		return "Generators"
	}
}

func aliasesFor(g Generator) []string {
	name := g.Name()
	aliases := []string{name}
	switch normalizeName(name) {
	case "person":
		aliases = append(aliases, "people")
	case "products":
		aliases = append(aliases, "product")
	case "names":
		aliases = append(aliases, "name")
	case "emails":
		aliases = append(aliases, "email")
	case "phonenumbers":
		aliases = append(aliases, "phone", "phones")
	case "addresses":
		aliases = append(aliases, "address")
	case "numbers":
		aliases = append(aliases, "number")
	case "identifiers":
		aliases = append(aliases, "identifier", "uuid", "ids")
	case "network":
		aliases = append(aliases, "net", "internet", "ip")
	case "companies":
		aliases = append(aliases, "company")
	case "dates":
		aliases = append(aliases, "date")
	case "colors":
		aliases = append(aliases, "color")
	case "photos":
		aliases = append(aliases, "photo", "images", "image")
	case "loremipsum":
		aliases = append(aliases, "lorem", "text")
	}
	return aliases
}

func normalizeName(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			out = append(out, r+'a'-'A')
		case r >= 'a' && r <= 'z':
			out = append(out, r)
		case r >= '0' && r <= '9':
			out = append(out, r)
		}
	}
	return string(out)
}

package generator

import (
	"fmt"
	"math/rand"
)

type namedColor struct {
	name    string
	r, g, b int
}

// Real color names mapped to their actual RGB values
var namedColors = []namedColor{
	{"Red", 220, 38, 38},
	{"Blue", 59, 130, 246},
	{"Green", 34, 197, 94},
	{"Yellow", 234, 179, 8},
	{"Purple", 168, 85, 247},
	{"Orange", 249, 115, 22},
	{"Pink", 236, 72, 153},
	{"Cyan", 6, 182, 212},
	{"Magenta", 217, 70, 239},
	{"Lime", 132, 204, 22},
	{"Teal", 20, 184, 166},
	{"Indigo", 99, 102, 241},
	{"Violet", 139, 92, 246},
	{"Coral", 251, 113, 133},
	{"Salmon", 250, 128, 114},
	{"Gold", 234, 179, 8},
	{"Silver", 163, 163, 163},
	{"Crimson", 220, 38, 38},
	{"Navy", 30, 58, 138},
	{"Olive", 101, 116, 60},
	{"Maroon", 127, 29, 29},
	{"Aqua", 34, 211, 238},
	{"Turquoise", 45, 212, 191},
	{"Lavender", 196, 181, 253},
	{"Beige", 245, 245, 220},
	{"Ivory", 255, 255, 240},
	{"Khaki", 189, 183, 107},
	{"Plum", 168, 85, 247},
	{"Orchid", 218, 112, 214},
	{"Slate", 100, 116, 139},
	{"Charcoal", 54, 69, 79},
	{"Midnight", 25, 25, 112},
	{"Forest", 22, 101, 52},
	{"Sky", 56, 189, 248},
	{"Rose", 244, 63, 94},
	{"Amber", 245, 158, 11},
	{"Emerald", 16, 185, 129},
	{"Ruby", 190, 18, 60},
	{"Sapphire", 37, 99, 235},
	{"Pearl", 234, 234, 234},
}

type ColorGen struct{}

func (g *ColorGen) Name() string        { return "Colors" }
func (g *ColorGen) Description() string { return "Color names with matching hex and RGB values" }
func (g *ColorGen) Fields() []Field {
	return []Field{
		{Name: "name", Desc: "Color name"},
		{Name: "hex", Desc: "#RRGGBB hex code"},
		{Name: "rgb", Desc: "rgb(r, g, b)"},
		{Name: "r", Desc: "Red channel 0-255"},
		{Name: "g", Desc: "Green channel 0-255"},
		{Name: "b", Desc: "Blue channel 0-255"},
	}
}

func (g *ColorGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		// Pick a named color — name, hex, and RGB all match
		c := namedColors[rng.Intn(len(namedColors))]

		// Add slight random variation (±15) to make each record unique
		// while staying recognizably the named color
		r := clamp(c.r + rng.Intn(31) - 15)
		g := clamp(c.g + rng.Intn(31) - 15)
		b := clamp(c.b + rng.Intn(31) - 15)

		records[i] = map[string]any{
			"name": c.name,
			"hex":  fmt.Sprintf("#%02x%02x%02x", r, g, b),
			"rgb":  fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
			"r":    r,
			"g":    g,
			"b":    b,
		}
	}
	return records
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

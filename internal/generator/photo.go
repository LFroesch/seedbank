package generator

import (
	"fmt"
	"math/rand"
)

var photoCategories = []string{
	"nature", "city", "people", "food", "animals", "architecture",
	"technology", "travel", "business", "abstract", "sports", "fashion",
}

var dimensions = [][2]int{
	{200, 200}, {300, 300}, {400, 400}, {640, 480}, {800, 600},
	{1024, 768}, {1280, 720}, {1920, 1080}, {150, 150}, {500, 500},
}

type PhotoGen struct{}

func (g *PhotoGen) Name() string        { return "Photos" }
func (g *PhotoGen) Description() string { return "Placeholder image URLs and metadata" }
func (g *PhotoGen) Fields() []Field {
	return []Field{
		{Name: "url", Desc: "Placeholder image URL"},
		{Name: "width", Desc: "Image width in px"},
		{Name: "height", Desc: "Image height in px"},
		{Name: "category", Desc: "Image category"},
		{Name: "alt_text", Desc: "Alt text description"},
	}
}

func (g *PhotoGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		dim := dimensions[rng.Intn(len(dimensions))]
		cat := photoCategories[rng.Intn(len(photoCategories))]
		w, h := dim[0], dim[1]
		seed := rng.Intn(10000)

		// picsum.photos is reliable and deterministic with seed
		url := fmt.Sprintf("https://picsum.photos/seed/%s-%d/%d/%d", cat, seed, w, h)
		alt := fmt.Sprintf("A %s photo (%dx%d)", cat, w, h)

		records[i] = map[string]any{
			"url":      url,
			"width":    w,
			"height":   h,
			"category": cat,
			"alt_text": alt,
		}
	}
	return records
}

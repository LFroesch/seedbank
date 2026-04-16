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

func (g *PhotoGen) Name() string { return "Photos" }
func (g *PhotoGen) Description() string {
	return "Field source for placeholder image URLs and metadata (WIP: semantic accuracy is still rough)"
}
func (g *PhotoGen) Kind() Kind { return KindField }
func (g *PhotoGen) Fields() []Field {
	return []Field{
		{Name: "url", Desc: "Seeded placeholder image URL"},
		{Name: "width", Desc: "Image width in px"},
		{Name: "height", Desc: "Image height in px"},
		{Name: "aspect_ratio", Desc: "Width:height ratio for the placeholder"},
		{Name: "category", Desc: "Loose category hint used in the seed"},
		{Name: "seed", Desc: "Deterministic placeholder seed"},
		{Name: "alt_text", Desc: "WIP caption derived from placeholder metadata"},
	}
}

func (g *PhotoGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		dim := dimensions[rng.Intn(len(dimensions))]
		cat := photoCategories[rng.Intn(len(photoCategories))]
		w, h := dim[0], dim[1]
		seed := rng.Intn(10000)
		seedTag := fmt.Sprintf("%s-%d", cat, seed)

		// picsum.photos is reliable and deterministic with seed, but not semantically exact.
		url := fmt.Sprintf("https://picsum.photos/seed/%s/%d/%d", seedTag, w, h)
		alt := fmt.Sprintf("WIP placeholder: %s scene at %dx%d", cat, w, h)

		records[i] = map[string]any{
			"url":          url,
			"width":        w,
			"height":       h,
			"aspect_ratio": fmt.Sprintf("%d:%d", w, h),
			"category":     cat,
			"seed":         seedTag,
			"alt_text":     alt,
		}
	}
	return records
}

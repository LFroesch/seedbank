package generator

import (
	"fmt"
	"math/rand"
	"time"
)

type DateGen struct{}

func (g *DateGen) Name() string { return "Dates" }
func (g *DateGen) Description() string {
	return "Field source for internally consistent dates, timestamps, and ages"
}
func (g *DateGen) Kind() Kind { return KindField }
func (g *DateGen) Fields() []Field {
	return []Field{
		{Name: "date", Desc: "YYYY-MM-DD"},
		{Name: "datetime", Desc: "YYYY-MM-DD HH:MM:SS"},
		{Name: "timestamp", Desc: "Unix timestamp"},
		{Name: "time", Desc: "HH:MM:SS"},
		{Name: "year", Desc: "Year (from date)"},
		{Name: "age", Desc: "Age if date were a birthday"},
	}
}

func (g *DateGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		// Generate a single random timestamp — all fields derive from it
		minT := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		maxT := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		ts := minT + rng.Int63n(maxT-minT)
		t := time.Unix(ts, 0).UTC()

		// Age derived against fixed reference date for stable seeded output.
		age := ageAt(t, referenceDate())

		records[i] = map[string]any{
			"date":      t.Format("2006-01-02"),
			"datetime":  t.Format("2006-01-02 15:04:05"),
			"timestamp": ts,
			"time":      t.Format("15:04:05"),
			"year":      t.Year(),
			"age":       age,
		}
	}
	return records
}

// FormatDuration formats seconds into a human-readable string.
func FormatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	}
	return fmt.Sprintf("%dh %dm", seconds/3600, (seconds%3600)/60)
}

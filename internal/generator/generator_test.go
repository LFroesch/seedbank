package generator

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestAllGeneratorsRegistered(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatal("no generators registered")
	}
	// Should have at least the 14 known generators
	if len(Registry) < 14 {
		t.Errorf("expected at least 14 generators, got %d", len(Registry))
	}
}

func TestGeneratorFieldsNotEmpty(t *testing.T) {
	for _, gen := range Registry {
		t.Run(gen.Name(), func(t *testing.T) {
			fields := gen.Fields()
			if len(fields) == 0 {
				t.Error("generator has no fields")
			}
			for _, f := range fields {
				if f.Name == "" {
					t.Error("field has empty name")
				}
			}
		})
	}
}

func TestGeneratorDeterministic(t *testing.T) {
	seed := int64(42)
	count := 5

	for _, gen := range Registry {
		t.Run(gen.Name(), func(t *testing.T) {
			rng1 := rand.New(rand.NewSource(seed))
			rng2 := rand.New(rand.NewSource(seed))

			r1 := gen.Generate(count, rng1)
			r2 := gen.Generate(count, rng2)

			if len(r1) != len(r2) {
				t.Fatalf("different record counts: %d vs %d", len(r1), len(r2))
			}

			for i := range r1 {
				for k, v1 := range r1[i] {
					v2, ok := r2[i][k]
					if !ok {
						t.Errorf("record %d missing key %q in second run", i, k)
						continue
					}
					// Compare as strings to handle all types
					if toString(v1) != toString(v2) {
						t.Errorf("record %d field %q: %v != %v", i, k, v1, v2)
					}
				}
			}
		})
	}
}

func TestGeneratorRecordCount(t *testing.T) {
	counts := []int{0, 1, 5, 50}
	rng := rand.New(rand.NewSource(99))

	for _, gen := range Registry {
		for _, count := range counts {
			t.Run(gen.Name(), func(t *testing.T) {
				records := gen.Generate(count, rng)
				if len(records) != count {
					t.Errorf("expected %d records, got %d", count, len(records))
				}
			})
		}
	}
}

func TestGeneratorFieldsPresent(t *testing.T) {
	rng := rand.New(rand.NewSource(123))

	for _, gen := range Registry {
		t.Run(gen.Name(), func(t *testing.T) {
			records := gen.Generate(3, rng)
			expectedFields := gen.Fields()
			for i, rec := range records {
				for _, f := range expectedFields {
					if _, ok := rec[f.Name]; !ok {
						t.Errorf("record %d missing field %q", i, f.Name)
					}
				}
			}
		})
	}
}

func TestEmailLinkedFields(t *testing.T) {
	gen := &EmailGen{}
	rng := rand.New(rand.NewSource(42))
	records := gen.Generate(5, rng)

	for i, rec := range records {
		if rec["first_name"] == nil || rec["first_name"] == "" {
			t.Errorf("record %d: first_name is empty", i)
		}
		if rec["last_name"] == nil || rec["last_name"] == "" {
			t.Errorf("record %d: last_name is empty", i)
		}
	}
}

func TestMixGenerator(t *testing.T) {
	gens := []Generator{&NameGen{}, &EmailGen{}}
	fields := BuildMixFields(gens)
	mix := &MixGenerator{Gens: gens, Fields_: fields}

	rng := rand.New(rand.NewSource(42))
	records := mix.Generate(3, rng)

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	// Should have prefixed fields from both generators
	for i, rec := range records {
		if rec["names.first_name"] == nil {
			t.Errorf("record %d: missing names.first_name", i)
		}
		if rec["emails.email"] == nil {
			t.Errorf("record %d: missing emails.email", i)
		}
	}
}

func TestPhotoGeneratorIncludesPlaceholderMetadata(t *testing.T) {
	gen := &PhotoGen{}
	rng := rand.New(rand.NewSource(7))
	record := gen.Generate(1, rng)[0]

	for _, field := range []string{"url", "width", "height", "aspect_ratio", "category", "seed", "alt_text"} {
		if _, ok := record[field]; !ok {
			t.Fatalf("photo record missing field %q", field)
		}
	}
	if got := toString(record["alt_text"]); len(got) < len("WIP placeholder") || got[:len("WIP placeholder")] != "WIP placeholder" {
		t.Fatalf("alt_text = %q, want WIP placeholder prefix", got)
	}
}

func TestPersonAgeMatchesDOB(t *testing.T) {
	gen := &PersonGen{}
	rng := rand.New(rand.NewSource(42))
	records := gen.Generate(25, rng)

	for i, rec := range records {
		dob, ok := rec["dob"].(string)
		if !ok {
			t.Fatalf("record %d dob is not a string", i)
		}
		age, ok := rec["age"].(int)
		if !ok {
			t.Fatalf("record %d age is not an int", i)
		}
		dobTime, err := time.Parse("2006-01-02", dob)
		if err != nil {
			t.Fatalf("record %d invalid dob %q: %v", i, dob, err)
		}
		want := ageAt(dobTime, referenceDate())
		if age != want {
			t.Fatalf("record %d age mismatch: got %d want %d", i, age, want)
		}
	}
}

func TestDateAgeMatchesDate(t *testing.T) {
	gen := &DateGen{}
	rng := rand.New(rand.NewSource(42))
	records := gen.Generate(25, rng)

	for i, rec := range records {
		date, ok := rec["date"].(string)
		if !ok {
			t.Fatalf("record %d date is not a string", i)
		}
		age, ok := rec["age"].(int)
		if !ok {
			t.Fatalf("record %d age is not an int", i)
		}
		tm, err := time.Parse("2006-01-02", date)
		if err != nil {
			t.Fatalf("record %d invalid date %q: %v", i, date, err)
		}
		want := ageAt(tm, referenceDate())
		if age != want {
			t.Fatalf("record %d age mismatch: got %d want %d", i, age, want)
		}
	}
}

func toString(v any) string {
	return fmt.Sprintf("%v", v)
}

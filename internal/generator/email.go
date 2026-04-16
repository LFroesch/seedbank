package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

type EmailGen struct{}

func (g *EmailGen) Name() string { return "Emails" }
func (g *EmailGen) Description() string {
	return "Field source for email addresses derived from fake names"
}
func (g *EmailGen) Kind() Kind { return KindField }
func (g *EmailGen) Fields() []Field {
	return []Field{
		{Name: "email", Desc: "Full email address"},
		{Name: "username", Desc: "Local part before @"},
		{Name: "domain", Desc: "Domain after @"},
		{Name: "first_name", Desc: "Name used to build the email"},
		{Name: "last_name", Desc: "Name used to build the email"},
	}
}

func (g *EmailGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		first, last, _, _ := pickGendered(rng)
		domain := safePersonalDomains[rng.Intn(len(safePersonalDomains))]

		firstL := strings.ToLower(first)
		lastL := strings.ToLower(last)

		var username string
		switch rng.Intn(4) {
		case 0:
			username = firstL + "." + lastL
		case 1:
			username = firstL[:1] + lastL
		case 2:
			username = firstL + fmt.Sprintf("%d", rng.Intn(999))
		case 3:
			username = lastL + "." + firstL[:1]
		}

		records[i] = map[string]any{
			"email":      username + "@" + domain,
			"username":   username,
			"domain":     domain,
			"first_name": first,
			"last_name":  last,
		}
	}
	return records
}

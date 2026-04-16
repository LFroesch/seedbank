package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

type PersonGen struct{}

func (g *PersonGen) Name() string { return "Person" }
func (g *PersonGen) Description() string {
	return "Coherent person record: identity, contact info, address, and company all line up"
}
func (g *PersonGen) Kind() Kind { return KindRecord }
func (g *PersonGen) Fields() []Field {
	return []Field{
		{Name: "id", Desc: "UUID"},
		{Name: "prefix", Desc: "Gender-matched prefix"},
		{Name: "first_name", Desc: "First name"},
		{Name: "last_name", Desc: "Last name"},
		{Name: "full_name", Desc: "First + Last"},
		{Name: "username", Desc: "Derived from name"},
		{Name: "email", Desc: "Derived from name + domain"},
		{Name: "phone", Desc: "Phone with state-matched area code"},
		{Name: "street", Desc: "Street address"},
		{Name: "city", Desc: "City (matches state)"},
		{Name: "state", Desc: "State code"},
		{Name: "zip", Desc: "Zip code (matches city)"},
		{Name: "country", Desc: "Country"},
		{Name: "company", Desc: "Employer"},
		{Name: "job_title", Desc: "Role at company"},
		{Name: "department", Desc: "Department"},
		{Name: "work_email", Desc: "Corporate email from name+company"},
		{Name: "website", Desc: "Company website"},
		{Name: "dob", Desc: "Date of birth"},
		{Name: "age", Desc: "Age derived from DOB"},
		{Name: "avatar", Desc: "Avatar URL from initials"},
	}
}

func (g *PersonGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		// Name — gender-coherent prefix
		first, last, prefix, _ := pickGendered(rng)
		firstL := strings.ToLower(first)
		lastL := strings.ToLower(last)

		// Username derived from name
		var username string
		switch rng.Intn(3) {
		case 0:
			username = firstL + "." + lastL
		case 1:
			username = firstL[:1] + lastL
		case 2:
			username = firstL + itoa(rng.Intn(99))
		}

		// Personal email — real providers are fine, the name is fake
		personalDomain := safePersonalDomains[rng.Intn(len(safePersonalDomains))]
		email := username + "@" + personalDomain

		// Location — coherent state/city/zip/area code
		state, city, zip, areaCode := pickLocation(rng)
		street := pickStreet(rng)

		// Phone — area code matches their state
		mid := 200 + rng.Intn(800)
		last4 := rng.Intn(10000)
		phone := fmt.Sprintf("(%s) %03d-%04d", areaCode, mid, last4)

		// Company
		compPrefix := companyPrefixes[rng.Intn(len(companyPrefixes))]
		compSuffix := companySuffixes[rng.Intn(len(companySuffixes))]
		company := compPrefix + " " + compSuffix
		dept := departments[rng.Intn(len(departments))]
		title := jobTitles[rng.Intn(len(jobTitles))]
		// Work domain uses .example.com to avoid collisions with real companies
		compSlug := strings.ToLower(strings.ReplaceAll(compPrefix, " ", ""))
		compDomain := compSlug + ".example.com"
		website := "https://www." + compDomain
		workEmail := firstL + "." + lastL + "@" + compDomain

		// DOB and age derived against a fixed reference date for stable seeded output.
		ref := referenceDate()
		minDob := ref.AddDate(-75, 0, 0)
		maxDob := ref.AddDate(-18, 0, 0)
		daySpan := int(maxDob.Sub(minDob).Hours() / 24)
		dobTime := minDob.AddDate(0, 0, rng.Intn(daySpan+1))
		dob := dobTime.Format("2006-01-02")
		age := ageAt(dobTime, ref)

		// Avatar from initials
		avatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s+%s&size=200&background=random&bold=true", first, last)

		records[i] = map[string]any{
			"id":         genUUID(rng),
			"prefix":     prefix,
			"first_name": first,
			"last_name":  last,
			"full_name":  first + " " + last,
			"username":   username,
			"email":      email,
			"phone":      phone,
			"street":     street,
			"city":       city,
			"state":      state,
			"zip":        zip,
			"country":    "United States",
			"company":    company,
			"job_title":  title,
			"department": dept,
			"work_email": workEmail,
			"website":    website,
			"dob":        dob,
			"age":        age,
			"avatar":     avatar,
		}
	}
	return records
}

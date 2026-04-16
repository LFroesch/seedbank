package generator

import (
	"math/rand"
	"strings"
)

var companyPrefixes = []string{
	"Apex", "Nova", "Vertex", "Quantum", "Stellar", "Zenith", "Prime", "Atlas",
	"Nexus", "Cipher", "Omega", "Pulse", "Core", "Summit", "Flux", "Vortex",
	"Catalyst", "Horizon", "Synapse", "Radiant", "Titan", "Echo", "Surge", "Onyx",
	"Beacon", "Forge", "Harbor", "Crest", "Ember", "Lumen", "Drift", "Pinnacle",
	"Cobalt", "Granite", "Aether", "Meridian", "Stratos", "Nimbus", "Prism", "Helix",
}

var companySuffixes = []string{
	"Systems", "Solutions", "Technologies", "Labs", "Industries", "Corp", "Inc",
	"Group", "Digital", "Analytics", "Software", "Dynamics", "Networks", "Partners",
	"Ventures", "Global", "Consulting", "Services", "Media", "Creative", "Studio",
	"Robotics", "Biotech", "Health", "Finance", "Capital", "Works", "Logic",
}

var departments = []string{
	"Engineering", "Marketing", "Sales", "Finance", "HR", "Operations",
	"Product", "Design", "Legal", "Support", "Research", "QA",
	"Data Science", "Security", "Infrastructure", "Customer Success",
}

var jobTitles = []string{
	"Software Engineer", "Product Manager", "Data Analyst", "UX Designer",
	"DevOps Engineer", "Marketing Manager", "Sales Rep", "QA Engineer",
	"Frontend Developer", "Backend Developer", "Full Stack Developer",
	"System Administrator", "Project Manager", "Business Analyst",
	"Technical Writer", "Security Engineer", "Data Scientist", "CTO",
	"VP of Engineering", "Team Lead", "Staff Engineer", "Principal Engineer",
	"Account Executive", "Solutions Architect", "Engineering Manager",
	"Director of Operations", "Chief of Staff", "Recruiter",
}

type CompanyGen struct{}

func (g *CompanyGen) Name() string { return "Companies" }
func (g *CompanyGen) Description() string {
	return "Field source for company, department, and job title data"
}
func (g *CompanyGen) Kind() Kind { return KindField }
func (g *CompanyGen) Fields() []Field {
	return []Field{
		{Name: "company", Desc: "Company name"},
		{Name: "department", Desc: "Department name"},
		{Name: "job_title", Desc: "Job title"},
		{Name: "website", Desc: "Company website URL"},
		{Name: "employee_count", Desc: "Number of employees"},
	}
}

func (g *CompanyGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		prefix := companyPrefixes[rng.Intn(len(companyPrefixes))]
		suffix := companySuffixes[rng.Intn(len(companySuffixes))]
		company := prefix + " " + suffix
		dept := departments[rng.Intn(len(departments))]
		title := jobTitles[rng.Intn(len(jobTitles))]
		// Use .example.com (RFC 2606 reserved) to avoid real domain collisions
		slug := strings.ToLower(strings.ReplaceAll(prefix, " ", ""))
		website := "https://www." + slug + ".example.com"
		empCount := 10 + rng.Intn(9990)

		records[i] = map[string]any{
			"company":        company,
			"department":     dept,
			"job_title":      title,
			"website":        website,
			"employee_count": empCount,
		}
	}
	return records
}

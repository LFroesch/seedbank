package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

var urlSchemes = []string{"https", "http"}
var hostLabels = []string{
	"api", "app", "cdn", "assets", "media", "auth", "admin", "checkout",
	"portal", "search", "files", "data", "static", "edge", "status",
}
var safeTLDs = []string{"example.com", "example.net", "example.org"}
var pathParts = []string{
	"v1", "v2", "users", "products", "orders", "billing", "settings", "health",
	"events", "feed", "items", "sessions", "accounts",
}
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/537.36 Chrome/126.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 Version/18.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Android 15; Mobile) AppleWebKit/537.36 Chrome/126.0.0.0 Mobile Safari/537.36",
	"curl/8.9.1",
}

type NetworkGen struct{}

func (g *NetworkGen) Name() string { return "Network" }
func (g *NetworkGen) Description() string {
	return "Field source for IP, MAC, hostname, URL, and user-agent data"
}
func (g *NetworkGen) Kind() Kind { return KindField }
func (g *NetworkGen) Fields() []Field {
	return []Field{
		{Name: "ipv4", Desc: "IPv4 address"},
		{Name: "ipv6", Desc: "IPv6 address"},
		{Name: "mac_address", Desc: "MAC address"},
		{Name: "hostname", Desc: "Safe hostname under RFC2606 domains"},
		{Name: "url", Desc: "URL built from generated host/path"},
		{Name: "user_agent", Desc: "Common browser or client user-agent"},
	}
}

func (g *NetworkGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		host := fmt.Sprintf("%s.%s", randomLabel(rng), safeTLDs[rng.Intn(len(safeTLDs))])
		scheme := urlSchemes[rng.Intn(len(urlSchemes))]
		path := fmt.Sprintf("/%s/%s", pathParts[rng.Intn(len(pathParts))], pathParts[rng.Intn(len(pathParts))])
		records[i] = map[string]any{
			"ipv4":        genIPv4(rng),
			"ipv6":        genIPv6(rng),
			"mac_address": genMAC(rng),
			"hostname":    host,
			"url":         fmt.Sprintf("%s://%s%s", scheme, host, path),
			"user_agent":  userAgents[rng.Intn(len(userAgents))],
		}
	}
	return records
}

func genIPv4(rng *rand.Rand) string {
	first := 11 + rng.Intn(212) // 11..222 avoids common reserved first octets
	if first == 127 {
		first = 126
	}
	return fmt.Sprintf("%d.%d.%d.%d", first, rng.Intn(256), rng.Intn(256), 1+rng.Intn(254))
}

func genIPv6(rng *rand.Rand) string {
	parts := make([]string, 8)
	for i := range parts {
		parts[i] = fmt.Sprintf("%x", rng.Intn(0x10000))
	}
	return strings.Join(parts, ":")
}

func genMAC(rng *rand.Rand) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		rng.Intn(256), rng.Intn(256), rng.Intn(256),
		rng.Intn(256), rng.Intn(256), rng.Intn(256),
	)
}

func randomLabel(rng *rand.Rand) string {
	base := hostLabels[rng.Intn(len(hostLabels))]
	if rng.Intn(2) == 0 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, 10+rng.Intn(90))
}

package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

type productCategory struct {
	name     string
	items    []string
	minPrice float64
	maxPrice float64
}

var productCategories = []productCategory{
	{"Electronics", []string{"Wireless Headphones", "Bluetooth Speaker", "USB-C Hub", "Mechanical Keyboard", "Webcam", "Monitor Stand", "Laptop Stand", "Power Bank", "Smart Watch", "Tablet Case"}, 19.99, 299.99},
	{"Clothing", []string{"Cotton T-Shirt", "Denim Jacket", "Running Shoes", "Wool Sweater", "Baseball Cap", "Leather Belt", "Linen Pants", "Flannel Shirt", "Rain Jacket", "Hoodie"}, 9.99, 149.99},
	{"Home", []string{"Scented Candle", "Throw Pillow", "Wall Clock", "Picture Frame", "Desk Lamp", "Plant Pot", "Cutting Board", "Coffee Mug", "Bath Towel Set", "Bookshelf"}, 7.99, 89.99},
	{"Books", []string{"Programming Guide", "Science Fiction Novel", "Cookbook", "Biography", "Self-Help Book", "History Book", "Art Book", "Travel Guide", "Mystery Novel", "Poetry Collection"}, 4.99, 49.99},
	{"Sports", []string{"Yoga Mat", "Resistance Bands", "Water Bottle", "Jump Rope", "Foam Roller", "Gym Bag", "Tennis Racket", "Basketball", "Cycling Gloves", "Running Armband"}, 9.99, 79.99},
	{"Food", []string{"Organic Coffee Beans", "Dark Chocolate Bar", "Olive Oil", "Hot Sauce", "Granola Mix", "Dried Mango", "Protein Bars", "Green Tea", "Honey Jar", "Trail Mix"}, 3.99, 34.99},
}

var productAdjectives = []string{
	"Premium", "Classic", "Pro", "Ultra", "Essential", "Deluxe", "Eco",
	"Compact", "Original", "Signature", "Vintage", "Modern", "Elite", "Basic",
}

// Varied description templates — {0}=adjective, {1}=item, {2}=category
var descTemplates = []string{
	"The %s %s is built for everyday %s use. Trusted by thousands of customers.",
	"Our %s %s combines quality and value for the %s enthusiast. A top-rated favorite.",
	"Introducing the %s %s — designed with care for your %s needs. Satisfaction guaranteed.",
	"Discover the %s %s, a standout in our %s collection. Built to last.",
	"Meet the %s %s. Engineered for the modern %s lifestyle. Free returns included.",
	"The %s %s delivers on both form and function for %s lovers. Highly recommended.",
}

type ProductGen struct{}

func (g *ProductGen) Name() string        { return "Products (Linked)" }
func (g *ProductGen) Description() string { return "Coherent products: name, SKU, category, price, rating — all connected" }
func (g *ProductGen) Fields() []Field {
	return []Field{
		{Name: "id", Desc: "Product UUID"},
		{Name: "sku", Desc: "SKU derived from category+name"},
		{Name: "name", Desc: "Product name"},
		{Name: "category", Desc: "Product category"},
		{Name: "price", Desc: "Price appropriate to category"},
		{Name: "currency", Desc: "Formatted price string"},
		{Name: "rating", Desc: "1.0-5.0 rating"},
		{Name: "reviews", Desc: "Review count (correlated with rating)"},
		{Name: "in_stock", Desc: "Stock boolean"},
		{Name: "stock_qty", Desc: "Quantity (0 when out of stock)"},
		{Name: "description", Desc: "Short description"},
		{Name: "image", Desc: "Deterministic product image"},
	}
}

func (g *ProductGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		// Pick a category — this drives everything
		cat := productCategories[rng.Intn(len(productCategories))]
		item := cat.items[rng.Intn(len(cat.items))]
		adj := productAdjectives[rng.Intn(len(productAdjectives))]
		name := adj + " " + item

		// Price within category range
		priceRange := cat.maxPrice - cat.minPrice
		price := cat.minPrice + rng.Float64()*priceRange
		price = float64(int(price*100)) / 100

		// SKU from category prefix + item initials + number
		catPrefix := strings.ToUpper(cat.name[:3])
		itemInitials := ""
		for _, word := range strings.Fields(item) {
			if len(word) > 0 {
				itemInitials += strings.ToUpper(word[:1])
			}
		}
		sku := fmt.Sprintf("%s-%s-%04d", catPrefix, itemInitials, rng.Intn(10000))

		// Rating 1.0-5.0, weighted toward 3-5
		rating := 1.0 + rng.Float64()*4.0
		rating = float64(int(rating*10)) / 10
		if rng.Intn(3) > 0 {
			rating = 2.5 + rng.Float64()*2.5
			rating = float64(int(rating*10)) / 10
		}

		// Reviews correlated with rating
		baseReviews := int(rating * 20)
		reviews := baseReviews + rng.Intn(baseReviews+1)

		// Stock
		inStock := rng.Intn(10) > 1
		stockQty := 0
		if inStock {
			stockQty = 1 + rng.Intn(500)
		}

		// Description — varied templates, no repeated adjective
		tmpl := descTemplates[rng.Intn(len(descTemplates))]
		desc := fmt.Sprintf(tmpl, adj, item, strings.ToLower(cat.name))

		// Image: picsum.photos with deterministic seed from SKU (always loads, always consistent)
		image := fmt.Sprintf("https://picsum.photos/seed/%s/400/400", sku)

		records[i] = map[string]any{
			"id":          genUUID(rng),
			"sku":         sku,
			"name":        name,
			"category":    cat.name,
			"price":       price,
			"currency":    fmt.Sprintf("$%.2f", price),
			"rating":      rating,
			"reviews":     reviews,
			"in_stock":    inStock,
			"stock_qty":   stockQty,
			"description": desc,
			"image":       image,
		}
	}
	return records
}

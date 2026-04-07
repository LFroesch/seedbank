package generator

import (
	"math/rand"
	"strings"
)

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum", "at", "vero", "eos",
	"accusamus", "iusto", "odio", "dignissimos", "ducimus", "blanditiis",
	"praesentium", "voluptatum", "deleniti", "atque", "corrupti", "quos", "dolores",
	"quas", "molestias", "excepturi", "obcaecati", "cupiditate", "provident",
	"similique", "ab", "illo", "inventore", "veritatis", "quasi", "architecto",
	"beatae", "vitae", "dicta", "explicabo", "nemo", "ipsam", "voluptatem", "quia",
	"voluptas", "aspernatur", "aut", "odit", "fugit", "consequuntur", "magni",
}

type LoremGen struct{}

func (g *LoremGen) Name() string        { return "Lorem Ipsum" }
func (g *LoremGen) Description() string  { return "Paragraphs, sentences, and words of placeholder text" }
func (g *LoremGen) Fields() []Field {
	return []Field{
		{Name: "word", Desc: "Single random word"},
		{Name: "sentence", Desc: "5-12 word sentence"},
		{Name: "paragraph", Desc: "3-6 sentence paragraph"},
		{Name: "title", Desc: "2-5 word capitalized title"},
	}
}

func (g *LoremGen) Generate(count int, rng *rand.Rand) []map[string]any {
	records := make([]map[string]any, count)
	for i := range records {
		records[i] = map[string]any{
			"word":      loremWords[rng.Intn(len(loremWords))],
			"sentence":  genSentence(rng, 5+rng.Intn(8)),
			"paragraph": genParagraph(rng, 3+rng.Intn(4)),
			"title":     genTitle(rng, 2+rng.Intn(4)),
		}
	}
	return records
}

func genSentence(rng *rand.Rand, wordCount int) string {
	words := make([]string, wordCount)
	for i := range words {
		words[i] = loremWords[rng.Intn(len(loremWords))]
	}
	s := strings.Join(words, " ")
	// Capitalize first letter
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s + "."
}

func genParagraph(rng *rand.Rand, sentenceCount int) string {
	sentences := make([]string, sentenceCount)
	for i := range sentences {
		sentences[i] = genSentence(rng, 5+rng.Intn(8))
	}
	return strings.Join(sentences, " ")
}

func genTitle(rng *rand.Rand, wordCount int) string {
	words := make([]string, wordCount)
	for i := range words {
		w := loremWords[rng.Intn(len(loremWords))]
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

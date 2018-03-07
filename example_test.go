// This example demonstrates an usage of suggest.Service
package suggest_test

import (
	"fmt"
	"github.com/alldroll/suggest"
)

// This example demonstrates how to use this package.
func Example() {
	// here we define our alphabet for given collection of words
	alphabet := suggest.NewCompositeAlphabet([]suggest.Alphabet{
		suggest.NewEnglishAlphabet(),
		suggest.NewSimpleAlphabet([]rune{'$'}), // pad wrap
	})

	// we create InMemoryDictionary. Here we can use anything we want,
	// for example SqlDictionary
	collection := []string{
		"Nissan March",
		"Nissan Juke",
		"Nissan Maxima",
		"Nissan Murano",
		"Nissan Note",
		"Toyota Mark II",
		"Toyota Corolla",
		"Toyota Corona",
	}

	dictionary := suggest.NewInMemoryDictionary(collection)

	// create IndexConfig with ngramSize, alphabet, wrap and pad
	wrap, pad := "$", "$"
	ngramSize := 3
	conf, err := suggest.NewIndexConfig(ngramSize, dictionary, alphabet, wrap, pad)
	if err != nil {
		panic(err)
	}

	service := suggest.NewService()
	service.AddRunTimeIndex("cars", conf)

	topK := 5
	sim := 0.4
	query := "niss ma"
	searchConf, err := suggest.NewSearchConfig(query, topK, suggest.CosineMetric(), sim)
	if err != nil {
		panic(err)
	}

	result, err := service.Suggest("cars", searchConf)
	if err != nil {
		panic(err)
	}

	values := make([]string, 0, len(result))
	for _, item := range result {
		values = append(values, item.Value)
	}

	fmt.Println(values)
	// Output: [Nissan Maxima Nissan March]
}

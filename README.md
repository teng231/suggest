# suggest

Tool for Top-k Approximate String Matching.

[![Go Report Card](https://goreportcard.com/badge/github.com/alldroll/suggest)](https://goreportcard.com/report/github.com/alldroll/suggest)
[![GoDoc](https://godoc.org/github.com/Lazin/go-ngram?status.png)](https://godoc.org/github.com/alldroll/suggest)

## Usage

```go
// here we define our alphabet for given collection of words
alphabet := suggest.NewCompositeAlphabet([]suggest.Alphabet{
    suggest.NewEnglishAlphabet(),
    suggest.NewSimpleAlphabet([]rune{'$'}), // pad wrap
})

// create IndexConfig with ngramSize, alphabet, wrap and pad
wrap, pad := "$", "$"
ngramSize := 3
conf, err := suggest.NewIndexConfig(ngramSize, alphabet, wrap, pad)
if err != nil {
    panic(err)
}

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

service := suggest.NewService()
service.AddDictionary("cars", dictionary, conf)

topK := 5
sim := 0.5
query := "niss ma"
searchConf, err := suggest.NewSearchConfig(query, topK, suggest.CosineMetric(), sim)
if err != nil {
    panic(err)
}

result := service.Suggest("cars", searchConf)
values := make([]string, 0, len(result))
for _, item := range result {
    values = append(values, item.Value)
}

fmt.Println(values)
// Output: [Nissan Maxima Nissan March]
```

## Demo
see https://tranquil-journey-12522.herokuapp.com/ as complete example
or run it localy by `go run cmd/web/main.go`

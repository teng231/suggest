package suggest

import (
	"fmt"

	"github.com/teng231/suggest/pkg/analysis"
	"github.com/teng231/suggest/pkg/dictionary"
	"github.com/teng231/suggest/pkg/index"
	"github.com/teng231/suggest/pkg/store"
)

// Index builds a search index by using the given config and the dictionary
// and persists it the directory
func Index(
	directory store.Directory,
	dict dictionary.Dictionary,
	config index.WriterConfig,
	tokenizer analysis.Tokenizer,
) error {
	encoder, err := index.NewEncoder()

	if err != nil {
		return fmt.Errorf("failed to create Encoder: %w", err)
	}

	indexWriter := index.NewIndexWriter(
		directory,
		config,
		encoder,
	)

	err = dict.Iterate(func(key dictionary.Key, value dictionary.Value) error {
		return indexWriter.AddDocument(key, tokenizer.Tokenize(value))
	})

	if err != nil {
		return err
	}

	if err = indexWriter.Commit(); err != nil {
		return err
	}

	return nil
}

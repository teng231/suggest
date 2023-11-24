package dep

import (
	"fmt"
	"github.com/teng231/suggest/pkg/dictionary"
	"github.com/teng231/suggest/pkg/lm"
	"github.com/teng231/suggest/pkg/spellchecker"
	"github.com/teng231/suggest/pkg/store"
	"github.com/teng231/suggest/pkg/suggest"
)

// BuildSpellChecker builds spellchecker for the provided config and indexDescription
func BuildSpellChecker(config *lm.Config, indexDescription suggest.IndexDescription) (*spellchecker.SpellChecker, error) {
	directory, err := store.NewFSDirectory(config.GetOutputPath())

	if err != nil {
		return nil, fmt.Errorf("failed to create a fs directory: %w", err)
	}

	languageModel, err := lm.RetrieveLMFromBinary(directory, config)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a lm model from binary format: %w", err)
	}

	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		return nil, fmt.Errorf("failed to open a cdb dictionary: %w", err)
	}

	// create runtime search index builder
	builder, err := suggest.NewRAMBuilder(dict, indexDescription)

	if err != nil {
		return nil, fmt.Errorf("failed to create a ngram index: %w", err)
	}

	index, err := builder.Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build a ngram index: %w", err)
	}

	return spellchecker.New(
		index,
		languageModel,
		lm.NewTokenizer(config.GetWordsAlphabet()),
		dict,
	), nil
}

package suggest

import (
	"github.com/alldroll/suggest/index"
	"github.com/alldroll/suggest/list_merger"
)

const (
	defaultPad       = "$"
	defaultWrap      = "$"
	defaultNGramSize = 3
)

// Builder
type Builder interface {
	Build() NGramIndex
}

// runTimeBuilderImpl implements Builder interface
type runTimeBuilderImpl struct {
	config *IndexConfig
}

// NewRunTimeBuilder returns new instance of runTimeBuilderImpl
func NewRunTimeBuilder(config *IndexConfig) Builder {
	return &runTimeBuilderImpl{
		config: config,
	}
}

func (b *runTimeBuilderImpl) Build() NGramIndex {
	conf := b.config
	cleaner := index.NewCleaner(conf.alphabet.Chars(), conf.pad, conf.wrap)
	generator := index.NewGenerator(conf.nGramSize, conf.alphabet)
	indexer := index.NewIndexer(
		conf.nGramSize,
		generator,
		cleaner,
	)

	indices := indexer.Index(conf.dictionary)
	indicesBuilder := index.NewInMemoryInvertedIndexIndicesBuilder(indices)

	return NewNGramIndex(
		cleaner,
		generator,
		indicesBuilder.Build(),
		&list_merger.CPMerge{},
	)
}

// builderImpl implements Builder interface
type builderImpl struct {
	description IndexDescription
}

// NewBuilder works with already indexed data
func NewBuilder(description IndexDescription) Builder {
	return &builderImpl{
		description: description,
	}
}

func (b *builderImpl) Build() NGramIndex {
	desc := b.description
	alphabet := desc.CreateAlphabet()

	cleaner := index.NewCleaner(alphabet.Chars(), desc.Pad, desc.Wrap)
	generator := index.NewGenerator(desc.NGramSize, alphabet)

	indicesBuilder := index.NewOnDiscInvertedIndexIndicesBuilder(
		desc.GetHeaderFile(),
		desc.GetDocumentListFile(),
	)

	return NewNGramIndex(
		cleaner,
		generator,
		indicesBuilder.Build(),
		&list_merger.CPMerge{},
	)
}

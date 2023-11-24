package index

import (
	"encoding/gob"
	"fmt"
	"runtime"

	"github.com/teng231/suggest/pkg/store"
)

// Reader is an entity, providing access to a search index
type Reader struct {
	directory store.Directory
	config    WriterConfig
}

// NewIndexReader returns a new instance of a search index reader
func NewIndexReader(
	directory store.Directory,
	config WriterConfig,
) *Reader {
	return &Reader{
		directory: directory,
		config:    config,
	}
}

// Read reads a inverted index indices from the given directory
func (ir *Reader) Read() (InvertedIndexIndices, error) {
	// read header
	header, err := ir.readHeader()

	if err != nil {
		return nil, err
	}

	documentReader, err := ir.directory.OpenInput(ir.config.DocumentListFileName)

	if err != nil {
		return nil, fmt.Errorf("failed to open document list: %w", err)
	}

	index, err := ir.createInvertedIndexIndices(header, documentReader)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inverted index: %w", err)
	}

	runtime.SetFinalizer(index, func(d interface{}) {
		documentReader.Close()
	})

	return index, nil
}

// readHeader reads an index header from the given directory
func (ir *Reader) readHeader() (*header, error) {
	headerReader, err := ir.directory.OpenInput(ir.config.HeaderFileName)

	if err != nil {
		return nil, fmt.Errorf("failed to open header: %w", err)
	}

	header := &header{}
	decoder := gob.NewDecoder(headerReader)

	if err = decoder.Decode(header); err != nil {
		return nil, fmt.Errorf("failed to retrieve header: %w", err)
	}

	if header.Version != IndexVersion {
		return nil, fmt.Errorf("index version mismatch, expected %s version", IndexVersion)
	}

	if err = headerReader.Close(); err != nil {
		return nil, fmt.Errorf("failed to close header file: %w", err)
	}

	return header, nil
}

// createInvertedIndexIndices creates new instance of InvertedIndexIndices from the given header
func (ir *Reader) createInvertedIndexIndices(header *header, documentReader store.Input) (InvertedIndexIndices, error) {
	// create inverted index structure slice
	indices := make([]InvertedIndex, int(header.Indices))
	invertedIndexStructureIndices := make([]invertedIndexStructure, len(indices))

	// here we create list of invertedIndexStructure
	for _, description := range header.Terms {
		if description.PostingListBytesSize == 0 {
			invertedIndexStructureIndices[description.Indice] = nil
			continue
		}

		if invertedIndexStructureIndices[description.Indice] == nil {
			invertedIndexStructureIndices[description.Indice] = make(invertedIndexStructure)
		}

		invertedIndexStructureIndices[description.Indice][description.Term] = struct {
			size     uint32
			position uint32
			length   uint32
		}{
			size:     description.PostingListBytesSize,
			position: description.PostingListPosition,
			length:   description.PostingListLen,
		}
	}

	// create NewInvertedIndex for given indice
	for i, invertedIndexStructure := range invertedIndexStructureIndices {
		if invertedIndexStructure == nil {
			indices[i] = nil
		} else {
			indices[i] = NewInvertedIndex(documentReader, invertedIndexStructure)
		}
	}

	return NewInvertedIndexIndices(indices), nil
}

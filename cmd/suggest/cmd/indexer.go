package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/teng231/suggest/pkg/store"

	"github.com/spf13/cobra"

	"github.com/teng231/suggest/pkg/dictionary"
	"github.com/teng231/suggest/pkg/suggest"
)

var (
	dict string
	host string
)

func init() {
	indexCmd.Flags().StringVarP(&dict, "dict", "d", "", "reindex certain dict")
	indexCmd.Flags().StringVarP(&host, "host", "", "", "host to send reindex request")

	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "indexer -c [config file]",
	Short: "builds indexes and dictionaries",
	Long:  `builds indexes and dictionaries and send signal to reload suggest-service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("indexer: ")
		log.SetFlags(0)

		configs, err := readConfigs()

		if err != nil {
			return err
		}

		reIndexed := false
		totalStart := time.Now()

		for _, config := range configs {
			if dict != "" && dict != config.Name {
				continue
			}

			if err := indexJob(config); err != nil {
				return err
			}

			reIndexed = true
		}

		if !reIndexed {
			log.Printf("There were not any reindex job")
			return nil
		}

		log.Printf("Total time spent %s", time.Since(totalStart).String())

		if pidPath != "" {
			if err := tryToSendReindexSignal(); err != nil {
				return err
			}
		}

		if host != "" {
			if err := tryToSendReindexRequest(); err != nil {
				return err
			}
		}

		return nil
	},
}

// readConfigs retrieves a list of index descriptions from the configPath
func readConfigs() ([]suggest.IndexDescription, error) {
	configs, err := suggest.ReadConfigs(configPath)

	if err != nil {
		return nil, fmt.Errorf("invalid config file format %w", err)
	}

	return configs, nil
}

// indexJob performs building a dictionary, a search index for the given index description
func indexJob(description suggest.IndexDescription) error {
	log.Printf("Start process '%s' config", description.Name)

	if description.Driver != suggest.DiscDriver {
		log.Printf("skip processing '%s', there is no disc configuration\n", description.Name)
		return nil
	}

	// create a cdb dictionary
	log.Printf("Building a dictionary...")
	start := time.Now()

	dict, err := buildDictionaryJob(description)

	if err != nil {
		return fmt.Errorf("failed to build a dictionary: %w", err)
	}

	log.Printf("Time spent %s", time.Since(start))

	// create a search index
	log.Printf("Creating a search index...")
	start = time.Now()

	directory, err := store.NewFSDirectory(description.GetIndexPath())

	if err != nil {
		return fmt.Errorf("failed to create a directory: %w", err)
	}

	if err = suggest.Index(directory, dict, description.GetWriterConfig(), description.GetIndexTokenizer()); err != nil {
		return err
	}

	log.Printf("Time spent %s", time.Since(start))
	log.Printf("End process\n\n")

	return nil
}

// buildDictionary builds a persistent dictionary
func buildDictionaryJob(config suggest.IndexDescription) (dictionary.Dictionary, error) {
	dictReader, err := newDictionaryReader(config)

	if err != nil {
		return nil, err
	}

	dict, err := dictionary.BuildCDBDictionary(dictReader, config.GetDictionaryFile())

	if err != nil {
		return nil, err
	}

	return dict, nil
}

// dictionaryReader is an adapter, that implements dictionary.Iterable for bufio.Scanner
type dictionaryReader struct {
	lineScanner *bufio.Scanner
}

// Iterate iterates through each line of the corresponding dictionary
func (dr *dictionaryReader) Iterate(iterator dictionary.Iterator) error {
	docID := dictionary.Key(0)

	for dr.lineScanner.Scan() {
		if err := iterator(docID, dr.lineScanner.Text()); err != nil {
			return err
		}

		docID++
	}

	return dr.lineScanner.Err()
}

// newDictionaryReader creates an adapter to Iterable interface, that scans all lines
// from the SourcePath and creates pairs of <DocID, Value>
func newDictionaryReader(config suggest.IndexDescription) (dictionary.Iterable, error) {
	f, err := os.Open(config.GetSourcePath())

	if err != nil {
		return nil, fmt.Errorf("could not open a source file %w", err)
	}

	scanner := bufio.NewScanner(f)

	return &dictionaryReader{
		lineScanner: scanner,
	}, nil
}

// tryToSendReindexSignal sends a SIGHUP signal to the pid
func tryToSendReindexSignal() error {
	d, err := ioutil.ReadFile(pidPath)

	if err != nil {
		return fmt.Errorf("error parsing pid from %s: %w", pidPath, err)
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))

	if err != nil {
		return fmt.Errorf("error parsing pid from %s: %w", pidPath, err)
	}

	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		return fmt.Errorf("fail to send reindex signal to %d, %w", pid, err)
	}

	return nil
}

// tryToSendReindexRequest sends a http request with a reindex purpose
func tryToSendReindexRequest() error {
	resp, err := http.Post(host, "text/plain", nil)

	if err != nil {
		return fmt.Errorf("fail to send reindex request to %s, %w", host, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("fail to read response body %w", err)
	}

	if string(body) != "OK" {
		return fmt.Errorf("something goes wrong with reindex request, %s", body)
	}

	return nil
}

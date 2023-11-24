package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/teng231/suggest/internal/http"
	"github.com/teng231/suggest/internal/spellchecker/dep"
	"github.com/teng231/suggest/pkg/lm"
	"github.com/teng231/suggest/pkg/suggest"
)

// App is our application
type App struct {
	config AppConfig
}

// AppConfig is an application config
type AppConfig struct {
	Port             string
	ConfigPath       string
	PidPath          string
	IndexDescription suggest.IndexDescription
}

// NewApp creates new instance of App for the given config
func NewApp(config AppConfig) App {
	return App{
		config: config,
	}
}

// Run starts the application
// performs http requests handling
func (a App) Run() error {
	config, err := lm.ReadConfig(a.config.ConfigPath)

	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	spellchecker, err := dep.BuildSpellChecker(config, a.config.IndexDescription)

	if err != nil {
		return err
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		a.listenToSystemSignals(cancelFn)
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", (&homeHandler{}).handle).Methods("GET")
	r.HandleFunc("/predict/{query}/", (&predictHandler{spellchecker}).handle).Methods("GET")

	corsHeaders := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET"})

	handler := handlers.LoggingHandler(os.Stdout, r)
	handler = handlers.CORS(corsHeaders, corsMethods)(handler)
	httpServer := http.NewServer(handler, "0.0.0.0:"+a.config.Port)

	return httpServer.Run(ctx)
}

// listenToSystemSignals handles OS signals
func (a App) listenToSystemSignals(cancelFn context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Println("Interrupt signal..")
			cancelFn()
		}
	}
}

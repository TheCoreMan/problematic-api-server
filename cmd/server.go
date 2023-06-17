package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/thecoreman/problematic-api-server/config"
	server "github.com/thecoreman/problematic-api-server/server/go"

	"github.com/gorilla/handlers"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultServerReadTimeout = 5 * time.Second
)

func main() {
	mainLogger := bootstrap()

	cachingAPIService := server.NewCachingApiService(mainLogger)
	cachingAPIController := server.NewCachingApiController(cachingAPIService)

	errorsAPIService := server.NewErrorsApiService()
	errorsAPIController := server.NewErrorsApiController(errorsAPIService)

	rateLimitingAPIService := server.NewRateLimitingApiService(mainLogger)
	rateLimitingAPIController := server.NewRateLimitingApiController(rateLimitingAPIService)

	router := server.NewRouter(cachingAPIController, errorsAPIController, rateLimitingAPIController)

	// Add CORS allowed origins
	allowedOrigins := viper.GetString(config.AllowedOriginsConfigKey)
	if allowedOrigins != "" {
		splitAllowedOrigins := strings.Split(allowedOrigins, ",")
		mainLogger.Info().
			Str("allowedOrigins", allowedOrigins).
			Msg("adding allowed origins")
		for _, origin := range splitAllowedOrigins {
			if origin == "*" {
				mainLogger.Warn().
					Str("allowedOrigins", allowedOrigins).
					Msg("Insecure! Allowed origins contains wildcard, CORS will allow all origins")
			}
		}
		router.Use(handlers.CORS(handlers.AllowedOrigins(splitAllowedOrigins)))
	} else {
		mainLogger.Info().
			Msg("CORS: no allowed origins configured")
	}

	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	serverAddress := getServerAddressFromConfig()

	server := &http.Server{
		Addr:              serverAddress,
		Handler:           router,
		ReadHeaderTimeout: defaultServerReadTimeout,
	}

	mainLogger.Info().
		Str("serverAddress", serverAddress).
		Msg("starting server 🚀")
	defer log.Info().Msg("server stopped 🛑")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// bootstrap sets up everything the server needs to run. For now
// it does configuration and logging.
//
// Anything that should happen on boot (versus runtime) should happen here.
func bootstrap() zerolog.Logger {
	config.BootstrapConfig()

	configuredLogLevel := viper.GetString(config.LogLevelConfigKey)
	level, err := zerolog.ParseLevel(configuredLogLevel)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("configuredLogLevel", configuredLogLevel).
			Msg("failed to parse log level")
		panic(err)
	}
	zerolog.SetGlobalLevel(level)

	if viper.GetString(config.LogFormatConfigKey) == "text" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return log.Logger.With().Str("component", "main").Logger()
}

func getServerAddressFromConfig() string {
	return fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
}

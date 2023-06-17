package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/thecoreman/problematic-api-server/config"
	"github.com/thecoreman/problematic-api-server/ratelimiter"
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
	// Add a mux Router that serves the /swagger directory on /swagger
	// This is a hack to get the swagger UI to work with the generated
	// server code.
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger/"))))
	mainLogger.Info().Str("swaggerPath", "/swagger/").Msg("serving swagger UI")

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

	ratelimiter := initializeRateLimiter(mainLogger)
	router.Use(ratelimiter.Middleware)

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
	serveErr := server.ListenAndServe()
	if serveErr != nil {
		panic(serveErr)
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

func initializeRateLimiter(mainLogger zerolog.Logger) ratelimiter.RateLimiter {
	ipLimitRateWindow, parseErr := time.ParseDuration(
		viper.GetString(config.RateLimitIPWindowConfigKey))
	if parseErr != nil {
		panic(parseErr)
	}
	accountLimitRateWindow, parseErr := time.ParseDuration(
		viper.GetString(config.RateLimitAccountWindowConfigKey))
	if parseErr != nil {
		panic(parseErr)
	}
	backoffLimitRateWindow, parseErr := time.ParseDuration(
		viper.GetString(config.RateLimitExponentialBackoffWindowConfigKey))
	if parseErr != nil {
		panic(parseErr)
	}
	backoffLimitFactor := viper.GetInt(config.RateLimitExponentialBackoffFactorConfigKey)
	ratelimiter := ratelimiter.NewRateLimiter(mainLogger, ratelimiter.Config{
		IPRateLimitWindow:          ipLimitRateWindow,
		AccountRateLimitWindow:     accountLimitRateWindow,
		BackoffRateLimitWindow:     backoffLimitRateWindow,
		BackoffRateLimitMultiplier: backoffLimitFactor,
	})
	return ratelimiter
}

func getServerAddressFromConfig() string {
	return fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
}

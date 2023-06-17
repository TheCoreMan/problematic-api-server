package config

import "github.com/spf13/viper"

const (
	envPrefix                   = "PROBLEMATIC"
	HostConfigKey               = "HOST"
	hostConfigDefaultValue      = "0.0.0.0"
	PortConfigKey               = "PORT"
	portConfigDefaultValue      = "4578"
	LogFormatConfigKey          = "LOG_FORMAT"
	logFormatConfigDefaultValue = "text"
	LogLevelConfigKey           = "LOG_LEVEL"
	logLevelConfigDefaultValue  = "debug"
	AllowedOriginsConfigKey     = "ALLOWED_ORIGINS"
	allowedOriginsDefaultValue  = "*"
	BooksDirectoryConfigKey     = "BOOKS_DIRECTORY"
	booksDirectoryDefaultValue  = "./books"
	//nolint:gosec // No hardcoded secrets here
	RateLimitIPWindowConfigKey                    = "RATE_LIMIT_IP_WINDOW"
	rateLimitIPWindowDefaultValue                 = "500ms"
	RateLimitAccountWindowConfigKey               = "RATE_LIMIT_ACCOUNT_WINDOW"
	rateLimitAccountWindowDefaultValue            = "500ms"
	RateLimitExponentialBackoffWindowConfigKey    = "RATE_LIMIT_EXPONENTIAL_BACKOFF_WINDOW"
	rateLimitExponentialBackoffDefaultValue       = "500ms"
	RateLimitExponentialBackoffFactorConfigKey    = "RATE_LIMIT_EXPONENTIAL_BACKOFF_FACTOR"
	rateLimitExponentialBackoffFactorDefaultValue = 2
)

func BootstrapConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	viper.SetDefault(HostConfigKey, hostConfigDefaultValue)
	viper.SetDefault(PortConfigKey, portConfigDefaultValue)
	viper.SetDefault(LogLevelConfigKey, logLevelConfigDefaultValue)
	viper.SetDefault(LogFormatConfigKey, logFormatConfigDefaultValue)
	viper.SetDefault(AllowedOriginsConfigKey, allowedOriginsDefaultValue)
	viper.SetDefault(BooksDirectoryConfigKey, booksDirectoryDefaultValue)
	viper.SetDefault(RateLimitIPWindowConfigKey, rateLimitIPWindowDefaultValue)
	viper.SetDefault(RateLimitAccountWindowConfigKey, rateLimitAccountWindowDefaultValue)
	viper.SetDefault(RateLimitExponentialBackoffWindowConfigKey, rateLimitExponentialBackoffDefaultValue)
	viper.SetDefault(RateLimitExponentialBackoffFactorConfigKey, rateLimitExponentialBackoffFactorDefaultValue)
}

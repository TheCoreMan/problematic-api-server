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
}

package config

var ServerConfigValues ServerConfig

// Model that links to config.yml file
type ServerConfig struct {
	APIDelegations struct {
		Endpoint                     string `yaml:"endpoint"`
		PollPeriodInSeconds          uint   `yaml:"pollPeriodInSeconds"`
		DelayLocalTimestampInSeconds uint   `yaml:"delayLocalTimestampInSeconds"`
	} `yaml:"api-delegations"`
	Database struct {
		Connections int    `yaml:"connections" env:"DB_CONNECTIONS" env-description:"Total number of database connections"`
		Name        string `yaml:"name" env:"DB_NAME" env-description:"Database name"`
		Host        string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
		Password    string `yaml:"pass"  env:"DB_PASSWORD" env-description:"db password"`
		Port        string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
		Username    string `yaml:"user"  env:"DB_USERNAME" env-description:"db username"`
		Timezone    string `yaml:"timezone" env:"DB_TIMEZONE" env-description:"Database timezone"`
	} `yaml:"database"`
	Server struct {
		APIPath            string   `yaml:"apiPath"  env:"API_PATH" env-description:"API base path"`
		APIVersion         string   `yaml:"apiVersion"  env:"API_VERSION" env-description:"API Version"`
		CORSAllowedClients []string `yaml:"corsAllowedClients" env:"CORS_ALLOWED_CLIENTS"  env-description:"List of allowed CORS Clients"`
		Environment        string   `yaml:"environment" env:"SERVER_ENVIRONMENT"  env-description:"server environment"`

		Host     string `yaml:"host"  env:"SERVER_HOST" env-description:"server host"`
		Port     string `yaml:"port" env:"SERVER_PORT"  env-description:"server port"`
		Protocol string `yaml:"protocol" env:"SERVER_PROTOCOL"  env-description:"server protocol"`
	} `yaml:"server"`
}

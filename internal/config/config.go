package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Sentry  SentryConfig  `toml:"sentry"`
	Servers ServersConfig `toml:"servers"`
	Clients ClientsConfig `toml:"clients"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

func (c GlobalConfig) IsProd() bool {
	return c.Env == "prod"
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"omitempty,url"`
}

type ServersConfig struct {
	Debug  DebugServerConfig `toml:"debug"`
	Client APIServerConfig   `toml:"client"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type APIServerConfig struct {
	Addr           string               `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `toml:"allow_origins" validate:"required"`
	RequiredAccess RequiredAccessConfig `toml:"required_access"`
}

type RequiredAccessConfig struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type ClientsConfig struct {
	Keycloak KeycloakConfig `toml:"keycloak"`
}

type KeycloakConfig struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required,alphanum"`
	DebugMode    bool   `toml:"debug_mode"`
}

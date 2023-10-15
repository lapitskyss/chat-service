package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Clients Clients       `toml:"clients"`
	Sentry  SentryConfig  `toml:"sentry"`
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

type ServersConfig struct {
	Client ClientServerConfig `toml:"client"`
	Debug  DebugServerConfig  `toml:"debug"`
}

type ClientServerConfig struct {
	Addr           string                     `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string                   `toml:"allow_origins" validate:"min=1"`
	RequiredAccess ClientServerRequiredAccess `toml:"required_access"`
}

type ClientServerRequiredAccess struct {
	Resource string `toml:"resource"`
	Role     string `toml:"role"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type Clients struct {
	Keycloak KeycloakClient `toml:"keycloak"`
}

type KeycloakClient struct {
	BasePath     string `toml:"base_path" validate:"required"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required,alphanum"`
	DebugMode    bool   `toml:"debug_mode"`
}

type SentryConfig struct {
	Dsn string `toml:"dsn" validate:"omitempty,url"`
}

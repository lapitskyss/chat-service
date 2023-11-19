package config

import "time"

type Config struct {
	Global   GlobalConfig   `toml:"global"`
	Log      LogConfig      `toml:"log"`
	Sentry   SentryConfig   `toml:"sentry"`
	PSQL     PSQLConfig     `toml:"psql"`
	Servers  ServersConfig  `toml:"servers"`
	Clients  ClientsConfig  `toml:"clients"`
	Services ServicesConfig `toml:"services"`
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

type PSQLConfig struct {
	Address  string `toml:"address" validate:"required,hostname_port"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password" validate:"required"`
	Database string `toml:"database" validate:"required"`
	Debug    bool   `toml:"debug"`
}

type ServersConfig struct {
	Debug   DebugServerConfig `toml:"debug"`
	Client  APIServerConfig   `toml:"client"`
	Manager APIManagerConfig  `toml:"manager"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type APIServerConfig struct {
	Addr           string               `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string             `toml:"allow_origins" validate:"required"`
	RequiredAccess RequiredAccessConfig `toml:"required_access"`
}

type APIManagerConfig struct {
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

type ServicesConfig struct {
	MsgProducer MsgProducerConfig `toml:"msg_producer"`
	Outbox      OutboxConfig      `toml:"outbox"`
	ManagerLoad ManagerLoadConfig `toml:"manager_load"`
}

type MsgProducerConfig struct {
	Brokers    []string `toml:"brokers" validate:"min=1"`
	Topic      string   `toml:"topic" validate:"required"`
	BatchSize  int      `toml:"batch_size" validate:"min=1,max=1000"`
	EncryptKey string   `toml:"encrypt_key" validate:"omitempty,hexadecimal"`
}

type OutboxConfig struct {
	Workers    int           `toml:"workers" validate:"min=1,max=32"`
	IdleTime   time.Duration `toml:"idle_time" validate:"min=100ms,max=10s"`
	ReserveFor time.Duration `toml:"reserve_for" validate:"min=1s,max=10m"`
}

type ManagerLoadConfig struct {
	MaxProblemsAtTime int `toml:"max_problems_at_same_time" validate:"min=1,max=30"`
}

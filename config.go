package draken

import (
	"os"
	"strings"
	"time"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Environment Environment
	Debug       bool
	Server      ServerConfig
	Storage     StorageConfig
	Redis       RedisConfig
}

type ServerConfig struct {
	Hidden    bool
	Port      uint16
	Heartbeat HeartbeatConfig
	Security  bool
}

type HeartbeatConfig struct {
	Enabled  bool
	Endpoint string
}

type Environment uint8

const (
	EnvironmentLocal Environment = iota
	EnvironmentDev
	EnvironmentStaging
	EnvironmentProd
)

type StorageType uint8

const (
	StorageTypeLibsql StorageType = iota
	StorageTypeSqlite
	StorageTypePostgres
)

type StorageConfig struct {
	Enabled bool
	Type    StorageType
	DSN     string
}

type RedisConfig struct {
	Enabled bool
	DSN     string
}

func (d *Draken) setup() error {
	d.StartedAt = time.Now()
	if err := d.loadConfigFile(); err != nil {
		return err
	}

	d.setDebug()
	d.setEnvironment()
	d.setLoggerOpts()
	d.setServerConfig()
	d.setStorageConfig()
	d.setRedisConfig()

	return nil
}

func (d *Draken) loadConfigFile() error {
	envFileFound := true
	log.Debug().Msgf("Loading environment variables...")
	if err := godotenv.Load(".env"); err != nil {
		envFileFound = false
	}

	log.Debug().Msgf("Loading config file...")
	raw, err := os.ReadFile(".config/draken.yaml")
	if err != nil {
		return errorx.DataUnavailable.New("loading config file failed")
	}

	substituted := string(raw)
	if envFileFound {
		log.Debug().Msgf("Substituting environment variables...")
		substituted, err = envsubst.String(string(raw))
		if err != nil {
			return errorx.RejectedOperation.New("substituting env variables failed")
		}
	}

	log.Debug().Msgf("Setting configuration...")
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(substituted)); err != nil {
		return errorx.RejectedOperation.New("reading config by viper failed")
	}

	log.Debug().Msgf("Registered keys %s in the configuration.", strings.Join(viper.AllKeys(), ", "))
	log.Info().Msgf("Configuration loaded.")
	return nil
}

func (d *Draken) setLoggerOpts() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if d.Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if d.Config.Environment == EnvironmentLocal {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func (d *Draken) setEnvironment() {
	str := viper.GetString("draken.environment")
	env := EnvironmentLocal
	switch str {
	case "dev":
		env = EnvironmentDev
	case "staging":
		env = EnvironmentStaging
	case "prod":
		env = EnvironmentProd
	default:
		env = EnvironmentLocal
	}
	d.Config.Environment = env
}

func (d *Draken) setDebug() {
	d.Config.Debug = viper.GetBool("draken.debug")
}

func (d *Draken) setStorageConfig() {
	enabled := viper.GetBool("draken.storage.enabled")
	str := viper.GetString("draken.storage.type")
	dsn := ""
	storageType := StorageTypeSqlite
	switch str {
	case "libsql":
		storageType = StorageTypeLibsql
		dsn = viper.GetString("draken.storage.libsql.dsn")
	case "sqlite":
		storageType = StorageTypeSqlite
	case "postgres":
		storageType = StorageTypePostgres
		dsn = viper.GetString("draken.storage.postgres.dsn")
	default:
		storageType = StorageTypeSqlite
	}
	d.Config.Storage.Enabled = enabled
	d.Config.Storage.DSN = dsn
	d.Config.Storage.Type = storageType
}

func (d *Draken) setRedisConfig() {
	d.Config.Redis.Enabled = viper.GetBool("draken.redis.enabled")
	d.Config.Redis.DSN = viper.GetString("draken.redis.dsn")
}

func (d *Draken) setServerConfig() {
	d.Config.Server.Port = viper.GetUint16("draken.server.port")
	d.Config.Server.Hidden = viper.GetBool("draken.server.hidden")
	d.Config.Server.Security = viper.GetBool("draken.server.security")
	d.Config.Server.Heartbeat.Enabled = viper.GetBool("draken.server.heartbeat.enabled")
	d.Config.Server.Heartbeat.Endpoint = viper.GetString("draken.server.heartbeat.endpoint")
}

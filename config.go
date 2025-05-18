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

type Environment uint8

const (
	EnvironmentLocal = iota
	EnvironmentDev
	EnvironmentStaging
	EnvironmentProd
)

func (d *Draken) setup() error {
	d.StartedAt = time.Now()
	if err := d.loadConfig(); err != nil {
		return err
	}

	d.setDebug()
	d.setEnvironment()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if d.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return nil
}

func (d *Draken) loadConfig() error {
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
	d.Environment = Environment(env)

}

func (d *Draken) setDebug() {
	d.Debug = viper.GetBool("draken.debug")
}

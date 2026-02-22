package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const configName = "goblocks"
const filename = configName + ".yaml"

type Config struct {
	Http   Http
	Blocks Blocks
}

func defaultConfig() {
	viper.SetDefault("http.host", "0.0.0.0")
	viper.SetDefault("http.port", 8000)
}

type Http struct {
	Host string
	Port int
}

type Blocks string

func (h *Config) HttpHostAndPort() string {
	return fmt.Sprintf("%s:%d", h.Http.Host, h.Http.Port)
}

func NewConfig() *Config {
	return loadConfig()
}
func loadConfig() *Config {
	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(
		".", "_"))

	defaultConfig()

	err := viper.ReadInConfig()

	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		viper.SafeWriteConfig()
	} else if err != nil {
		panic(err)
	}

	viper.AutomaticEnv()

	C := &Config{}
	err = viper.Unmarshal(C)
	if err != nil {
		panic("unable to decode into struct," + err.Error())
	}

	return C

}

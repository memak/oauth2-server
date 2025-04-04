package config

import (
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("values")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./oauth2-helm") // Ensure this path is correct

	// allow ENV vars to override
	_ = godotenv.Load()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Only read YAML config if file exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("values.yaml not found, continuing without it.")
		} else {
			log.Fatalf("Failed to read config: %v", err)
		}
	}
}

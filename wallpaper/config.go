package wallpaper

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WallpaperManagerConfig struct {
	Provider string
	Unsplash UnsplashConfig
	Bing     BingConfig
}

type UnsplashConfig struct {
	Sources    []string
	SearchTags []string
}

type BingConfig struct {
	Market string
}

func LoadConfig() *WallpaperManagerConfig {
	var config WallpaperManagerConfig

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("provider", "unsplash")
	viper.SetDefault("unsplash.sources", []string{"featured"})
	viper.SetDefault("unsplash.searchTags", []string{"wallpaper"})
	viper.SetDefault("bing.market", "en-US")

	log.Debug("Loading config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Info("config.yaml missing, will use defaults")
		} else {
			// Config file was found but another error was produced
			log.Error("Error reading config file - ", err)
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Unable to decode into struct - ", err)
	}

	return &config
}

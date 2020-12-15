package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func InitViper(configLocation, filename, extension string) {
	CheckConfigFile(configLocation, filename+extension)
	// Set up viper config library
	viper.SetConfigName(filename)
	viper.AddConfigPath(configLocation)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func CheckConfigFile(path, filename string) {
	//Check to see if the config file exists, if not, download and save it
	fullPath := filepath.Join(path, filename)
	if !fileExists(fullPath) {
		CreateFile(fullPath)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Info().Msg("File not found...")
		return false
	}
	log.Info().Msg("File found...")
	return !info.IsDir()
}

func CreateFile(filename string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("There was an error while creating the file. %s", filename))
		return
	}
}

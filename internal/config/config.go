package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

var suffixes = [...]string{"B","KB","MB","GB","TB"}

func Init() {
	viper.SetConfigName("stream")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func BufferMaxSize() int {
	maxSizeString := viper.GetString("buffer.max_size")
	for _, suffix := range suffixes {
		if strings.Contains(maxSizeString, suffix) {
			s := strings.Replace(maxSizeString, suffix, "", 1)
			maxSize, err := strconv.Atoi(s)
			if err != nil {
				continue
			}
			switch suffix {
			case suffixes[0]:
				return maxSize
			case suffixes[1]:
				return maxSize * 1024
			case suffixes[2]:
				return maxSize * (1024 * 1024)
			case suffixes[3]:
				return maxSize * (1024 * 1024 * 1024)
			case suffixes[4]:
				return maxSize * (1024 * 1024 * 1024)
			}
		}
	}

	log.Printf("invalid size %s\n", maxSizeString)
	panic(maxSizeString)
}
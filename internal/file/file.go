package file

import (
	"os"

	"github.com/spf13/viper"
)

func GetFile() (*os.File, error) {
	file, err := os.Open(viper.GetString("data.path"))
	if err != nil {
		return nil, err
	}

	return file, nil
}

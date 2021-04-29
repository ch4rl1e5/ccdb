package file

import (
	"github.com/spf13/viper"
	"os"
)

func GetFile() (*os.File, error) {
	file, err := os.Open(viper.GetString("file.path"))
	if err != nil {
		return nil, err
	}

	return file, nil
}

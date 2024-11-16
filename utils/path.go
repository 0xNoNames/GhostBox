package utils

import (
	"os"
)

func Exist(filesPath string) bool {
	_, err := os.Stat(filesPath)
	return err == nil
}

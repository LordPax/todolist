package utils

import (
	"os"
)

const (
	DB_FILE = "tasks.db"
)

func FileExist(file string) bool {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

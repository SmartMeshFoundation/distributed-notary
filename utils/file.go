package utils

import (
	"os"
	"runtime"
)

// Exists returns true if directory||file exists
func Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetHomePath returns the user's $HOME directory
func GetHomePath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

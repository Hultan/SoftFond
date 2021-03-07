package tools

import (
	"log"
	"os"
	"path"
)

func GetResourcePath(resource string) string {
	exePath := getExePath()
	// Try exePath/resource
	tryPath := path.Join(exePath, resource)
	if fileExists(tryPath) {
		return tryPath
	}
	// Try exePath/../resource
	tryPath = path.Join(exePath, "..", resource)
	if fileExists(tryPath) {
		return tryPath
	}
	return resource
}

func fileExists(existsPath string) bool {
	if _, err := os.Stat(existsPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func getExePath() string {
	exePath, err := os.Executable()
	if err!=nil {
		log.Fatal(err)
	}
	return path.Dir(exePath)
}


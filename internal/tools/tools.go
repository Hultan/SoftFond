package tools

import (
	"os"
	"path"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetAssetsPath(fileName string) (string, error) {
	return getResourcePath("assets", fileName)
}

func GetConfigPath(fileName string) (string, error) {
	return getResourcePath("config", fileName)
}

func getResourcePath(searchPath, fileName string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := path.Dir(exePath)

	gladePath := path.Join(exeDir, fileName)
	if FileExists(gladePath) {
		return gladePath, nil
	}
	gladePath = path.Join(exeDir, searchPath, fileName)
	if FileExists(gladePath) {
		return gladePath, nil
	}
	gladePath = path.Join(exeDir, "..", searchPath, fileName)
	if FileExists(gladePath) {
		return gladePath, nil
	}
	return gladePath, nil
}

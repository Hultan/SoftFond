package tools

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func MultiplySpaces(num int) string {
	return MultiplyString(num, " ")
}

func MultiplyString(num int, text string) string {
	var result string
	for i := 0; i <= num; i++ {
		result += text
	}
	return result
}

func GetRateString(text string) string {
	text = text[GetFirstNumberPosition(text):]
	return strings.Replace(text,",",".",1)
}

func GetFirstNumberPosition(text string) int {
	for i:=len(text)-1;i>=0;i-- {
		ascii := int(text[i])
		if ascii >= 48 && ascii <= 57 || ascii==44 {
			continue
		}
		return i+1
	}
	return -1
}

// GetResourcePath : Gets the path to a resource file
func GetResourcePath(directory, file string) string {
	return path.Join(GetExecutablePath(), directory, file)
}

// GetExecutablePath : Returns the path of the executable
func GetExecutablePath() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(executable)
}

package service

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var validFileSuffix = map[string]bool{
	".txt":  true,
	".md":   true,
	".go":   true,
	".js":   true,
	".java": true,
	".cpp":  true,
	".php":  true,
	".py":   true,
	".sh":   true,
	".c":    true,
	".h":    true,
	".conf": true,
	".json": true,
	".xml":  true,
	".yml":  true,
	".yaml": true,
	".csv":  true,
	".ini":  true,
	".toml": true,
}

func ReadFileContent(path string) (string, error) {

	//get file extention name
	ext := strings.ToLower(filepath.Ext(path))

	if condition, ok := validFileSuffix[ext]; !ok || !condition {
		return "", errors.New("unsupported file extention")
	}

	bytes, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(bytes), nil

}

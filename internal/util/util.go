package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadData(path string) string {
	var dataStr string
	if FileExists(path) {
		data, err := os.ReadFile(path)
		CheckErr(err)
		dataStr = string(data)
	} else {
		dataStr = path
	}
	return strings.TrimSpace(dataStr)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func BaseName(path string) string {
	return filepath.Base(path)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

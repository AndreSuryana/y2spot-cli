package utils

import (
	"os"
	"path/filepath"
	"time"
)

func WriteErrorLog(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	dir := "logs"
	_ = os.MkdirAll(dir, os.ModePerm)

	filename := time.Now().Format("log-2006-01-02-15-04-05.txt")
	fullPath := filepath.Join(dir, filename)

	f, _ := os.Create(fullPath)
	defer f.Close()

	for _, line := range lines {
		f.WriteString(line + "\n")
	}

	return fullPath
}

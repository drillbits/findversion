package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func chomp(s string) string {
	return strings.TrimRight(s, " \t\r\n")
}

func fmtPath(format string, a ...interface{}) string {
	return filepath.Clean(fmt.Sprintf(format, a...))
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

func fatal(w io.Writer, format string, a ...interface{}) int {
	fmt.Fprintf(w, "go tool dist: %s\n", fmt.Sprintf(format, a...))
	return ExitCodeError
}

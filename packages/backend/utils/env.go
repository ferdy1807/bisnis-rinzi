package utils

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv membaca file .env dan mengatur environment variables jika belum diatur.
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err // file tidak ditemukan atau masalah permission
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Lewati baris kosong atau komentar
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Jika belum ada di OS env, kita set
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}

	return scanner.Err()
}

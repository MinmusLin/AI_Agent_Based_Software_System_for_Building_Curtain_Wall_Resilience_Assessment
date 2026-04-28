package configs

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	CoreApiAddr string
	CoreBizAddr string
}

func Load() Config {
	return Config{
		CoreApiAddr: env("ICW_CORE_API_ADDR", "127.0.0.1:8000"),
		CoreBizAddr: env("ICW_CORE_BIZ_ADDR", "127.0.0.1:8001"),
	}
}

func LoadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open .env file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || value == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

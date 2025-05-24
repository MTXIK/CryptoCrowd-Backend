package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Config - основная структура конфигурации приложения
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

// DatabaseConfig - конфигурация базы данных
type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

// ServerConfig - конфигурация HTTP-сервера
type ServerConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	ReadTimeout     int    `json:"read_timeout"`     // в секундах
	WriteTimeout    int    `json:"write_timeout"`    // в секундах
	IdleTimeout     int    `json:"idle_timeout"`     // в секундах
	ShutdownTimeout int    `json:"shutdown_timeout"` // в секундах
}

// Load загружает конфигурацию из JSON-файла
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла конфигурации: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка декодирования файла конфигурации: %w", err)
	}

	// Устанавливаем значения по умолчанию, если они не определены
	setDefaults(&cfg)

	return &cfg, nil
}

// setDefaults устанавливает значения по умолчанию для параметров, которые не были заданы
func setDefaults(cfg *Config) {
	// Значения по умолчанию для сервера
	if cfg.Server.Host == "" {
		cfg.Server.Host = "127.0.0.1"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "9000"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 // 30 секунд
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 // 30 секунд
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 10 // 10 секунд
	}
	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 10 // 10 секунд
	}
}

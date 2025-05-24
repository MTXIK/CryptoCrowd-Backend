package main

import (
	"context"
	"fmt"
	"github.com/CryptoCrowd/internal/config"
	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/logger"
	"github.com/CryptoCrowd/migrate"
)

const (
	configPath = "config.json"
)

func main() {
	ctx, mainCtxCancel := context.WithCancel(context.Background())
	defer mainCtxCancel()

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	err = logger.InitGlobalLogger(cfg)
	if err != nil {
		logger.Fatalf("ошибка инициализации логгера: %v", err)
	}

	logger.Info("Приложение запускается...")
	logger.Infof("Используется конфигурация из файла: %s", configPath)

	pool, err := initInfrastructure(ctx, cfg)
	if err != nil {
		logger.Fatalf("ошибка инициализации инфраструктуры: %v", err)
	}
	defer pool.Close()
}

// Инициализация инфраструктуры (миграции, подключение к БД)
func initInfrastructure(ctx context.Context, cfg *config.Config) (*db.Pool, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.Name)

	logger.Infof("Запуск миграций для БД: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	if err := migrate.Run(dbURL, "up"); err != nil {
		return nil, fmt.Errorf("ошибка при выполнении миграций: %v", err)
	}
	logger.Debug("Миграции успешно выполнены")

	logger.Debug("Подключение к базе данных...")
	pool, err := db.NewPool(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	logger.Info("Подключение к базе данных установлено")

	return pool, nil
}

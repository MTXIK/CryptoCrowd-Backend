package main

import (
	"context"
	"fmt"
	"github.com/CryptoCrowd/internal/config"
	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/handler"
	"github.com/CryptoCrowd/internal/logger"
	"github.com/CryptoCrowd/internal/repository"
	"github.com/CryptoCrowd/internal/router"
	"github.com/CryptoCrowd/internal/service"
	"github.com/CryptoCrowd/migrate"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	configPath = "config.json"
)

type repositories struct {
	accRepo  *repository.PostgresAccount
	projRepo *repository.PostgresProject
	invRepo  *repository.PostgresInvestment
}

type services struct {
	accService  *service.Account
	projService *service.Project
	invService  *service.Investment
}

type handlers struct {
	accHandler  *handler.AccountHandler
	projHandler *handler.ProjectHandler
	invHandler  *handler.InvestmentHandler
}

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

	repos := initRepositories(pool)
	logger.Debug("Репозитории успешно инициализированы")

	services := initServices(repos)
	logger.Debug("Сервисы успешно инициализированы")

	handlers := initHandlers(services)
	logger.Debug("Хендлеры успешно инициализированы")

	app := router.SetupRouter(handlers.accHandler, handlers.projHandler, handlers.invHandler)
	logger.Debug("Маршруты успешно настроены")

	serverShutdown := startServer(ctx, app, cfg.Server.Port)
	defer serverShutdown()

	waitForShutdownSignal()

	logger.Info("Приложение успешно завершено")
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

func initRepositories(pool *db.Pool) *repositories {
	return &repositories{
		accRepo:  repository.NewPostgresAccount(pool),
		projRepo: repository.NewPostgresProject(pool),
		invRepo:  repository.NewPostgresInvestment(pool),
	}
}

func initServices(repos *repositories) *services {
	return &services{
		accService:  service.NewAccount(repos.accRepo),
		projService: service.NewProject(repos.projRepo),
		invService:  service.NewInvestment(repos.invRepo, repos.projRepo),
	}
}

func initHandlers(services *services) *handlers {
	return &handlers{
		accHandler:  handler.NewAccountHandler(services.accService),
		projHandler: handler.NewProjectHandler(services.projService),
		invHandler:  handler.NewInvestmentHandler(services.invService),
	}
}

func startServer(ctx context.Context, app *fiber.App, port string) func() {
	go func() {
		logger.Infof("HTTP-сервер запущен на порту :%s", port)
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatalf("ошибка запуска HTTP-сервера: %v", err)
		}
	}()

	return func() {
		logger.Debug("Остановка HTTP-сервера...")
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			logger.Fatalf("ошибка при завершении HTTP-сервера: %v", err)
		}
		logger.Debug("HTTP-сервер остановлен")
	}
}

func waitForShutdownSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	logger.Infof("Получен сигнал завершения %v, выключаем сервер...", sig)
}

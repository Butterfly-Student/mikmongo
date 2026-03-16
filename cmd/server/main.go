package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mikmongo/internal/config"
	"mikmongo/internal/domain"
	"mikmongo/internal/handler"
	"mikmongo/internal/middleware"
	_ "mikmongo/internal/migration"
	"mikmongo/internal/queue"
	"mikmongo/internal/queue/consumer"
	"mikmongo/internal/repository"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/router"
	"mikmongo/internal/scheduler"
	"mikmongo/internal/seeder"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/logger"
	"mikmongo/pkg/rabbitmq"
	"mikmongo/pkg/redis"
	casbinpkg "mikmongo/internal/casbin"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	"go.uber.org/zap"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	logg := logger.New(cfg.App.Env)
	defer logg.Sync()

	logg.Info("Starting MikMongo server...",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Database
	db, err := gorm.Open(gormpg.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		logg.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate + seed (when AUTO_MIGRATE=true)
	if cfg.Seed.AutoMigrate {
		sqlDB, err := db.DB()
		if err != nil {
			logg.Fatal("Failed to get sql.DB for migration", zap.Error(err))
		}
		if err := goose.SetDialect("postgres"); err != nil {
			logg.Fatal("Failed to set goose dialect", zap.Error(err))
		}
		if err := goose.Up(sqlDB, "."); err != nil {
			logg.Fatal("Failed to run migrations", zap.Error(err))
		}
		logg.Info("Migrations completed")

		seedCfg := seeder.Config{
			AdminEmail:     cfg.Seed.AdminEmail,
			AdminPassword:  cfg.Seed.AdminPassword,
			AdminName:      cfg.Seed.AdminName,
			AdminPhone:     cfg.Seed.AdminPhone,
			RouterName:     cfg.Seed.RouterName,
			RouterAddress:  cfg.Seed.RouterAddress,
			RouterAPIPort:  cfg.Seed.RouterAPIPort,
			RouterUsername: cfg.Seed.RouterUsername,
			RouterPassword: cfg.Seed.RouterPassword,
			EncryptionKey:  cfg.JWT.Secret,
		}
		if err := seeder.New(sqlDB, seedCfg).Run(context.Background()); err != nil {
			logg.Fatal("Failed to run seeder", zap.Error(err))
		}
		logg.Info("Seed completed")
	}

	// Redis
	redisClient := redis.NewClient(redis.Options{
		Host:     cfg.Redis.Host,
		Port:     parseInt(cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	// RabbitMQ
	var rabbitClient *rabbitmq.Client
	rabbitClient, err = rabbitmq.NewClient(
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
	)
	if err != nil {
		logg.Error("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		defer rabbitClient.Close()
	}

	// JWT
	jwtExpiry, _ := time.ParseDuration(cfg.JWT.Expiry)
	jwtRefreshExpiry, _ := time.ParseDuration(cfg.JWT.RefreshExpiry)
	if jwtRefreshExpiry == 0 {
		jwtRefreshExpiry = 7 * 24 * time.Hour
	}
	jwtService := jwt.NewService(cfg.JWT.Secret, jwtExpiry, jwtRefreshExpiry)

	// Repositories
	pgRepo := postgres.NewRepository(db)
	repoRegistry := &repository.Registry{
		CustomerRepo:             pgRepo.CustomerRepo,
		InvoiceRepo:              pgRepo.InvoiceRepo,
		PaymentRepo:              pgRepo.PaymentRepo,
		RouterDeviceRepo:         pgRepo.RouterDeviceRepo,
		UserRepo:                 pgRepo.UserRepo,
		BandwidthProfileRepo:     pgRepo.BandwidthProfileRepo,
		SubscriptionRepo:         pgRepo.SubscriptionRepo,
		CustomerRegistrationRepo: pgRepo.CustomerRegistrationRepo,
		InvoiceItemRepo:          pgRepo.InvoiceItemRepo,
		PaymentAllocationRepo:    pgRepo.PaymentAllocationRepo,
		SystemSettingRepo:        pgRepo.SystemSettingRepo,
		SequenceCounterRepo:      pgRepo.SequenceCounterRepo,
		MessageTemplateRepo:      pgRepo.MessageTemplateRepo,
		AuditLogRepo:             pgRepo.AuditLogRepo,
	}

	// Domains
	domainRegistry := domain.NewRegistry(
		domain.NewCustomerDomain(),
		domain.NewBillingDomain(),
		domain.NewPaymentDomain(),
		domain.NewRouterDomain(),
		domain.NewSubscriptionDomain(),
		domain.NewRegistrationDomain(),
		domain.NewNotificationDomain(),
	)

	// Queue
	queueRegistry := queue.NewRegistry(rabbitClient)

	// Services
	encKey := cfg.JWT.Secret
	serviceRegistry := service.NewRegistry(
		repoRegistry,
		domainRegistry,
		jwtService,
		encKey,
		db,
		redisClient,
		logg.Logger,
	)

	// Inject services into queue consumers (callbacks, avoids import cycle)
	queueRegistry.BillingConsumer.SetHandler(func(ctx context.Context, subscriptionID uuid.UUID, period time.Time) error {
		_, err := serviceRegistry.Billing.GenerateInvoice(ctx, subscriptionID, period)
		return err
	})
	queueRegistry.SuspendConsumer.SetHandler(func(ctx context.Context, customerID uuid.UUID) error {
		return serviceRegistry.Customer.IsolateAllSubscriptions(ctx, customerID)
	})
	queueRegistry.NotificationConsumer.SetHandler(&consumer.NotificationHandler{
		SendWhatsApp: func(ctx context.Context, phone, message string) error {
			return serviceRegistry.Notification.SendViaWhatsApp(ctx, phone, message)
		},
		SendEmail: func(ctx context.Context, to, subject, body string) error {
			return serviceRegistry.Notification.SendViaEmail(ctx, to, subject, body)
		},
	})


	// Handlers
	handlerRegistry := handler.NewRegistry(serviceRegistry, repoRegistry.SystemSettingRepo, jwtService)


	// Casbin RBAC enforcer
	casbinEnforcer, err := casbinpkg.NewEnforcer(db)
	if err != nil {
		logg.Fatal("Failed to create Casbin enforcer", zap.Error(err))
	}

	// Middleware
	middlewareRegistry := middleware.NewRegistry(logg.Logger, jwtService, redisClient, casbinEnforcer)

	// Router
	r := router.New(handlerRegistry, middlewareRegistry)

	// Schedulers
	schedulerRegistry := scheduler.NewRegistry(serviceRegistry, queueRegistry)
	schedulerRegistry.Start()
	defer schedulerRegistry.Stop()

	// HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	go func() {
		logg.Info("Server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("Server forced to shutdown", zap.Error(err))
	}

	logg.Info("Server exited")
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go                    # Viper: load .env / yaml
│   │
│   ├── migration/                       # ✅ BARU: Goose Go-based migrations
│   │   ├── registry.go                  # Import semua migration file (init trick)
│   │   ├── 0001_customers.go
│   │   ├── 0002_packages.go
│   │   ├── 0003_invoices.go
│   │   ├── 0004_payments.go
│   │   └── 0005_router_devices.go
│   │
│   ├── model/                           # GORM model structs
│   │   ├── customer.go
│   │   ├── package.go
│   │   ├── invoice.go
│   │   ├── payment.go
│   │   └── router_device.go
│   │
│   ├── domain/                          # Business logic
│   │   ├── registry.go
│   │   ├── customer/
│   │   │   └── domain.go
│   │   ├── billing/
│   │   │   └── domain.go
│   │   ├── payment/
│   │   │   └── domain.go
│   │   └── router/
│   │       └── domain.go
│   │
│   ├── repository/                      # ✅ LEBIH SIMPEL dengan GORM
│   │   ├── interfaces.go                # Repository interface registry
│   │   ├── customer_repo.go             # Interface
│   │   ├── invoice_repo.go
│   │   ├── payment_repo.go
│   │   ├── router_device_repo.go
│   │   └── postgres/                   # GORM implementations
│   │       ├── registry.go              # NewRepository(db *gorm.DB)
│   │       ├── customer_repo.go
│   │       ├── invoice_repo.go
│   │       ├── payment_repo.go
│   │       └── router_device_repo.go
│   │
│   ├── service/
│   │   ├── registry.go
│   │   ├── customer_service.go
│   │   ├── billing_service.go
│   │   ├── payment_service.go
│   │   └── router_service.go           # Orchestrate DB + pkg/mikrotik
│   │
│   ├── handler/
│   │   ├── registry.go
│   │   ├── customer_handler.go
│   │   ├── billing_handler.go
│   │   ├── payment_handler.go
│   │   ├── router_handler.go
│   │   └── webhook_handler.go          # Midtrans webhook
│   │
│   ├── queue/                          # ✅ BARU: RabbitMQ consumers & producers
│   │   ├── registry.go                  # Setup exchange, queue, binding
│   │   ├── producer/
│   │   │   ├── billing_producer.go      # Publish: generate invoice event
│   │   │   ├── suspend_producer.go      # Publish: suspend customer event
│   │   │   └── notification_producer.go
│   │   └── consumer/
│   │       ├── billing_consumer.go      # Consume: proses invoice
│   │       ├── suspend_consumer.go      # Consume: eksekusi suspend ke Mikrotik
│   │       └── notification_consumer.go # Consume: kirim email/WA
│   │
│   ├── scheduler/
│   │   ├── registry.go
│   │   ├── billing_scheduler.go        # Cron → publish ke RabbitMQ
│   │   ├── suspend_scheduler.go        # Cron → publish ke RabbitMQ
│   │   └── sync_scheduler.go           # Cron → sync data Mikrotik ke DB
│   │
│   ├── middleware/
│   │   ├── auth.go                     # JWT validation
│   │   ├── logger.go                   # Zap request logger
│   │   └── ratelimit.go                # Redis-based rate limiter
│   │
│   └── router/
│       └── router.go                   # Gin route definitions
│
├── pkg/                                # Reusable libraries
│   │
│   ├── mikrotik/                       # RouterOS client (sudah dibahas)
│   │   ├── domain/
│   │   │   ├── ppp.go
│   │   │   ├── hotspot.go
│   │   │   ├── queue.go
│   │   │   ├── firewall.go
│   │   │   ├── interface.go
│   │   │   └── errors.go
│   │   ├── client/
│   │   │   ├── client.go
│   │   │   ├── manager.go                 # Multi-router connection pool
│   │   │   └── options.go
│   │   ├── ppp/
│   │   ├── hotspot/
│   │   ├── queue/
│   │   ├── firewall/
│   │   ├── monitor/
│   │   └── mikrotik.go                 # Facade
│   │
│   ├── redis/                          # ✅ BARU: Redis client wrapper
│   │   ├── client.go                   # Connect, options, health check
│   │   ├── cache.go                    # Get, Set, Del, TTL helpers
│   │   ├── session.go                  # JWT session management
│   │   └── ratelimit.go                # Sliding window rate limiter
│   │
│   ├── rabbitmq/                       # ✅ BARU: RabbitMQ client wrapper
│   │   ├── client.go                   # Connect, reconnect, channel pool
│   │   ├── publisher.go                # Publish message ke exchange
│   │   ├── subscriber.go               # Subscribe & consume queue
│   │   ├── options.go                  # ExchangeOptions, QueueOptions
│   │   └── errors.go
│   │
│   ├── logger/
│   │   └── logger.go                   # Zap setup (dev vs production mode)
│   ├── jwt/
│   │   └── jwt.go                      # Sign & verify token
│   ├── response/
│   │   └── response.go                 # Standar JSON response wrapper
│   ├── pagination/
│   │   └── pagination.go
│   └── validator/
│       └── validator.go
│
├── tests/
│   ├── mocks/
│   │   ├── repository/
│   │   └── service/
│   └── integration/
│
├── deployments/
│   ├── Dockerfile
│   ├── docker-compose.yml              # Postgres + Redis + RabbitMQ
│   └── nginx.conf
│
├── .env.example
├── go.mod
├── go.sum
└── Makefile
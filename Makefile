IMAGE_NAME=$(shell basename $(CURDIR)):latest
CONTAINER_NAME=$(shell basename $(CURDIR))_app
MODULE=$(shell head -1 go.mod | awk '{print $$2}')

# .PHONY declares targets that don't create files with the same name as the target
.PHONY: build http \
	model domain service handler repository scheduler migration \
	queue-producer queue-consumer \
	mikrotik-domain mikrotik-module \
	migrate-up migrate-down migrate-status migrate-reset \
	seed fresh \
	generate-mocks tidy lint test test-coverage test-integration run

# ─────────────────────────────────────────────────────────────────────────────
# DOCKER
# ─────────────────────────────────────────────────────────────────────────────

build:
	@if [ "$(BUILD)" = "true" ]; then \
		echo "[INFO] BUILD=true, force rebuilding Docker image $(IMAGE_NAME)..."; \
		docker build -t $(IMAGE_NAME) .; \
	elif ! docker image inspect $(IMAGE_NAME) > /dev/null 2>&1; then \
		echo "[INFO] Docker image $(IMAGE_NAME) not found. Building..."; \
		docker build -t $(IMAGE_NAME) .; \
	else \
		echo "[INFO] Docker image $(IMAGE_NAME) already exists. Skipping build."; \
	fi

http:
	$(MAKE) build BUILD=$(BUILD)
	@echo "[INFO] Running the application in HTTP server mode inside Docker."
	docker run --rm \
	  --name $(CONTAINER_NAME) \
	  --env-file .env \
	  -p 8000:8000 \
	  --network $(shell basename $(CURDIR))_default \
	  $(IMAGE_NAME) http

docker-build:
	docker build -f deployments/Dockerfile -t mikmongo:latest .

docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker-compose.yml down

# ─────────────────────────────────────────────────────────────────────────────
# MIGRATE  (Goose — Go-based migrations in internal/migration/)
# Requires : go install github.com/pressly/goose/v3/cmd/goose@latest
# DSN      : read from .env  →  DB_DSN="postgres://user:pass@host:5432/db?sslmode=disable"
# ─────────────────────────────────────────────────────────────────────────────

DB_DSN ?= $(shell grep -E '^DB_DSN=' .env 2>/dev/null | sed 's/^DB_DSN=//' | tr -d '\r')

migrate-up:
	@echo "[INFO] Running all pending migrations..."
	@goose -dir internal/migration postgres "$(DB_DSN)" up

migrate-down:
	@echo "[INFO] Rolling back last migration..."
	@goose -dir internal/migration postgres "$(DB_DSN)" down

migrate-reset:
	@echo "[WARN] Resetting ALL migrations (down to 0)..."
	@goose -dir internal/migration postgres "$(DB_DSN)" reset

migrate-status:
	@echo "[INFO] Migration status:"
	@goose -dir internal/migration postgres "$(DB_DSN)" status

seed:
	@echo "[INFO] Running migrations + seed..."
	@go run cmd/seed/main.go

fresh:
	@$(MAKE) migrate-reset
	@$(MAKE) seed

# ─────────────────────────────────────────────────────────────────────────────
# MODEL
# Usage : make model VAL=customer
# Output: internal/model/customer.go
# ─────────────────────────────────────────────────────────────────────────────

model:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make model VAL=customer"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	DST=internal/model/$${LOWER}.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	mkdir -p internal/model; \
	printf "package model\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"time\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"gorm.io/gorm\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "type $${PASCAL} struct {\n" >> $$DST; \
	printf "\tID        uint           \`json:\"id\"         gorm:\"primaryKey\"\`\n" >> $$DST; \
	printf "\t$${PASCAL}Input\n" >> $$DST; \
	printf "\tCreatedAt time.Time      \`json:\"created_at\"\`\n" >> $$DST; \
	printf "\tUpdatedAt time.Time      \`json:\"updated_at\"\`\n" >> $$DST; \
	printf "\tDeletedAt gorm.DeletedAt \`json:\"-\"          gorm:\"index\"\`\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "type $${PASCAL}Input struct {\n" >> $$DST; \
	printf "\t// TODO: add fields\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "type $${PASCAL}Filter struct {\n" >> $$DST; \
	printf "\tIDs []uint \`json:\"ids\"\`\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func ($${PASCAL}) TableName() string { return \"$${LOWER}s\" }\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func (f $${PASCAL}Filter) IsEmpty() bool {\n" >> $$DST; \
	printf "\treturn len(f.IDs) == 0\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created model file: $$DST"

# ─────────────────────────────────────────────────────────────────────────────
# DOMAIN
# Usage : make domain VAL=billing
# Output: internal/domain/billing/domain.go
#         internal/domain/registry.go  (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

domain:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make domain VAL=billing"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	DOMAIN_DIR=internal/domain/$$LOWER; \
	if [ -d "$$DOMAIN_DIR" ]; then \
		echo "[ERROR] Directory $$DOMAIN_DIR already exists."; \
		exit 1; \
	fi; \
	mkdir -p $$DOMAIN_DIR; \
	echo "[INFO] Created directory: $$DOMAIN_DIR"; \
	DOMAIN_FILE=$$DOMAIN_DIR/domain.go; \
	printf "// Package %s contains %s domain logic.\n" "$$LOWER" "$$LOWER" >> $$DOMAIN_FILE; \
	printf "package $${LOWER}\n" >> $$DOMAIN_FILE; \
	printf "\n" >> $$DOMAIN_FILE; \
	printf "// Domain represents $${LOWER} business logic.\n" >> $$DOMAIN_FILE; \
	printf "type Domain struct{}\n" >> $$DOMAIN_FILE; \
	printf "\n" >> $$DOMAIN_FILE; \
	printf "// NewDomain creates a new $${LOWER} domain.\n" >> $$DOMAIN_FILE; \
	printf "func NewDomain() *Domain {\n" >> $$DOMAIN_FILE; \
	printf "\treturn &Domain{}\n" >> $$DOMAIN_FILE; \
	printf "}\n" >> $$DOMAIN_FILE; \
	echo "[INFO] Created domain file: $$DOMAIN_FILE"; \
	REGISTRY_FILE=internal/domain/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		mkdir -p internal/domain; \
		printf "package domain\n\nimport (\n)\n\n// Registry holds all domain instances.\ntype Registry struct {\n}\n\n// NewRegistry creates a new domain registry.\nfunc NewRegistry() *Registry {\n\treturn &Registry{}\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created domain registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "\"$(MODULE)/internal/domain/$$LOWER\"" "$$REGISTRY_FILE"; then \
		echo "[INFO] Import $$LOWER already exists"; \
	else \
		awk '/^import \($$/{print;print "\t\"$(MODULE)/internal/domain/'"$$LOWER"'\"";next}1' \
			"$$REGISTRY_FILE" > "$$REGISTRY_FILE.tmp" && mv "$$REGISTRY_FILE.tmp" "$$REGISTRY_FILE"; \
		echo "[INFO] Added import $$LOWER to domain registry"; \
	fi; \
	if grep -q "$${PASCAL} \*$${LOWER}.Domain" "$$REGISTRY_FILE"; then \
		echo "[INFO] Field $${PASCAL} already exists"; \
	else \
		awk '/^type Registry struct \{$$/{print;print "\t'"$${PASCAL}"' *'"$${LOWER}"'.Domain";next}1' \
			"$$REGISTRY_FILE" > "$$REGISTRY_FILE.tmp" && mv "$$REGISTRY_FILE.tmp" "$$REGISTRY_FILE"; \
		echo "[INFO] Added $${PASCAL} field to domain registry"; \
	fi; \
	if grep -q "func New$${PASCAL}Domain()" "$$REGISTRY_FILE"; then \
		echo "[INFO] New$${PASCAL}Domain() factory already exists"; \
	else \
		printf "\n// New$${PASCAL}Domain creates a new $${LOWER} domain.\nfunc New$${PASCAL}Domain() *$${LOWER}.Domain {\n\treturn $${LOWER}.NewDomain()\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Added New$${PASCAL}Domain() to domain registry"; \
	fi; \
	echo "[INFO] NOTE: wire $${PASCAL}: New$${PASCAL}Domain() in registry.go NewRegistry()"; \
	echo "[INFO] Domain generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# SERVICE
# Usage : make service VAL=billing
# Output: internal/service/billing_service.go
#         internal/service/registry.go  (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

service:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make service VAL=billing"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p internal/service; \
	DST=internal/service/$${LOWER}_service.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package service\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"$(MODULE)/internal/repository\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Service handles $${LOWER} business logic.\n" >> $$DST; \
	printf "type $${PASCAL}Service struct {\n" >> $$DST; \
	printf "\trepo repository.$${PASCAL}Repository\n" >> $$DST; \
	printf "\t// TODO: add other dependencies (domain, queue producer, etc.)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// New$${PASCAL}Service creates a new $${LOWER} service.\n" >> $$DST; \
	printf "func New$${PASCAL}Service(repo repository.$${PASCAL}Repository) *$${PASCAL}Service {\n" >> $$DST; \
	printf "\treturn &$${PASCAL}Service{repo: repo}\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created service file: $$DST"; \
	REGISTRY_FILE=internal/service/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "package service\n\n// Registry holds all service instances.\ntype Registry struct {\n}\n\n// NewRegistry creates a new service registry.\nfunc NewRegistry() *Registry {\n\treturn &Registry{}\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created service registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "$${PASCAL} \*$${PASCAL}Service" "$$REGISTRY_FILE"; then \
		echo "[INFO] Field $${PASCAL} already exists"; \
	else \
		awk '/^type Registry struct \{$$/{print;print "\t'"$${PASCAL}"' *'"$${PASCAL}"'Service";next}1' \
			"$$REGISTRY_FILE" > "$$REGISTRY_FILE.tmp" && mv "$$REGISTRY_FILE.tmp" "$$REGISTRY_FILE"; \
		echo "[INFO] Added $${PASCAL} field to service registry"; \
		echo "[INFO] NOTE: wire New$${PASCAL}Service(repo.$${PASCAL}Repo) in registry.go NewRegistry()"; \
	fi; \
	echo "[INFO] Service generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# HANDLER
# Usage : make handler VAL=billing
# Output: internal/handler/billing_handler.go
#         internal/handler/registry.go  (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

handler:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make handler VAL=billing"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p internal/handler; \
	DST=internal/handler/$${LOWER}_handler.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package handler\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"net/http\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"github.com/gin-gonic/gin\"\n" >> $$DST; \
	printf "\t\"$(MODULE)/internal/service\"\n" >> $$DST; \
	printf "\t\"$(MODULE)/pkg/response\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Handler handles HTTP requests for $${LOWER}.\n" >> $$DST; \
	printf "type $${PASCAL}Handler struct {\n" >> $$DST; \
	printf "\tservice *service.$${PASCAL}Service\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// New$${PASCAL}Handler creates a new $${LOWER} handler.\n" >> $$DST; \
	printf "func New$${PASCAL}Handler(s *service.$${PASCAL}Service) *$${PASCAL}Handler {\n" >> $$DST; \
	printf "\treturn &$${PASCAL}Handler{service: s}\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// List$${PASCAL} GET /$${LOWER}s\n" >> $$DST; \
	printf "func (h *$${PASCAL}Handler) List$${PASCAL}(c *gin.Context) {\n" >> $$DST; \
	printf "\tresponse.Success(c, http.StatusOK, \"ok\", nil)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Get$${PASCAL} GET /$${LOWER}s/:id\n" >> $$DST; \
	printf "func (h *$${PASCAL}Handler) Get$${PASCAL}(c *gin.Context) {\n" >> $$DST; \
	printf "\tresponse.Success(c, http.StatusOK, \"ok\", nil)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Create$${PASCAL} POST /$${LOWER}s\n" >> $$DST; \
	printf "func (h *$${PASCAL}Handler) Create$${PASCAL}(c *gin.Context) {\n" >> $$DST; \
	printf "\tresponse.Success(c, http.StatusCreated, \"created\", nil)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Update$${PASCAL} PUT /$${LOWER}s/:id\n" >> $$DST; \
	printf "func (h *$${PASCAL}Handler) Update$${PASCAL}(c *gin.Context) {\n" >> $$DST; \
	printf "\tresponse.Success(c, http.StatusOK, \"updated\", nil)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Delete$${PASCAL} DELETE /$${LOWER}s/:id\n" >> $$DST; \
	printf "func (h *$${PASCAL}Handler) Delete$${PASCAL}(c *gin.Context) {\n" >> $$DST; \
	printf "\tresponse.Success(c, http.StatusOK, \"deleted\", nil)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created handler file: $$DST"; \
	REGISTRY_FILE=internal/handler/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "package handler\n\nimport (\n\t\"$(MODULE)/internal/service\"\n)\n\n// Registry holds all handler instances.\ntype Registry struct {\n}\n\n// NewRegistry creates a new handler registry.\nfunc NewRegistry(services *service.Registry) *Registry {\n\treturn &Registry{}\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created handler registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "$${PASCAL} \*$${PASCAL}Handler" "$$REGISTRY_FILE"; then \
		echo "[INFO] Field $${PASCAL} already exists"; \
	else \
		awk '/^type Registry struct \{$$/{print;print "\t'"$${PASCAL}"' *'"$${PASCAL}"'Handler";next}1' \
			"$$REGISTRY_FILE" > "$$REGISTRY_FILE.tmp" && mv "$$REGISTRY_FILE.tmp" "$$REGISTRY_FILE"; \
		echo "[INFO] Added $${PASCAL} field to handler registry"; \
		echo "[INFO] NOTE: wire $${PASCAL}: New$${PASCAL}Handler(services.$${PASCAL}) in registry.go NewRegistry()"; \
	fi; \
	echo "[INFO] Handler generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# REPOSITORY
# Usage : make repository VAL=customer
# Output: internal/repository/customer_repo.go         (interface)
#         internal/repository/postgres/customer_repo.go (GORM impl)
#         internal/repository/registry.go               (auto-updated)
#         internal/repository/postgres/registry.go      (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

repository:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make repository VAL=customer"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
		CAMEL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_"} {printf "%s", $$1; for(i=2;i<=NF;i++) printf "%s", toupper(substr($$i,1,1)) substr($$i,2)} END{print ""}'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
		CAMEL=$$LOWER; \
	fi; \
	mkdir -p internal/repository internal/repository/postgres; \
	IFACE_FILE=internal/repository/$${LOWER}_repo.go; \
	if [ -f "$$IFACE_FILE" ]; then \
		echo "[ERROR] File $$IFACE_FILE already exists."; \
		exit 1; \
	fi; \
	printf "package repository\n" >> $$IFACE_FILE; \
	printf "\n" >> $$IFACE_FILE; \
	printf "import (\n" >> $$IFACE_FILE; \
	printf "\t\"context\"\n" >> $$IFACE_FILE; \
	printf "\n" >> $$IFACE_FILE; \
	printf "\t\"$(MODULE)/internal/model\"\n" >> $$IFACE_FILE; \
	printf ")\n" >> $$IFACE_FILE; \
	printf "\n" >> $$IFACE_FILE; \
	printf "//go:generate mockgen -source=$${LOWER}_repo.go -destination=./../../tests/mocks/repository/mock_$${LOWER}_repo.go\n" >> $$IFACE_FILE; \
	printf "// $${PASCAL}Repository defines data access for $${LOWER}.\n" >> $$IFACE_FILE; \
	printf "type $${PASCAL}Repository interface {\n" >> $$IFACE_FILE; \
	printf "\tCreate(ctx context.Context, v *model.$${PASCAL}) error\n" >> $$IFACE_FILE; \
	printf "\tGetByID(ctx context.Context, id uint) (*model.$${PASCAL}, error)\n" >> $$IFACE_FILE; \
	printf "\tList(ctx context.Context, limit, offset int) ([]model.$${PASCAL}, error)\n" >> $$IFACE_FILE; \
	printf "\tUpdate(ctx context.Context, v *model.$${PASCAL}) error\n" >> $$IFACE_FILE; \
	printf "\tDelete(ctx context.Context, id uint) error\n" >> $$IFACE_FILE; \
	printf "}\n" >> $$IFACE_FILE; \
	echo "[INFO] Created repository interface: $$IFACE_FILE"; \
	REGISTRY_FILE=internal/repository/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "// Package repository defines repository interfaces and registries.\npackage repository\n\n// Registry holds all repository interfaces.\ntype Registry struct {\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created repository registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "$${PASCAL}Repo $${PASCAL}Repository" "$$REGISTRY_FILE"; then \
		echo "[INFO] $${PASCAL}Repo already in repository registry"; \
	else \
		awk '/^type Registry struct \{$$/{print;print "\t'"$${PASCAL}"'Repo '"$${PASCAL}"'Repository";next}1' \
			"$$REGISTRY_FILE" > "$$REGISTRY_FILE.tmp" && mv "$$REGISTRY_FILE.tmp" "$$REGISTRY_FILE"; \
		echo "[INFO] Added $${PASCAL}Repo to repository registry"; \
	fi; \
	PG_FILE=internal/repository/postgres/$${LOWER}_repo.go; \
	if [ -f "$$PG_FILE" ]; then \
		echo "[ERROR] File $$PG_FILE already exists."; \
		exit 1; \
	fi; \
	printf "package postgres\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "import (\n" >> $$PG_FILE; \
	printf "\t\"context\"\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "\t\"gorm.io/gorm\"\n" >> $$PG_FILE; \
	printf "\t\"$(MODULE)/internal/model\"\n" >> $$PG_FILE; \
	printf "\t\"$(MODULE)/internal/repository\"\n" >> $$PG_FILE; \
	printf ")\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "type $${CAMEL}Repository struct {\n" >> $$PG_FILE; \
	printf "\tdb *gorm.DB\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "// New$${PASCAL}Repository creates a new $${LOWER} repository.\n" >> $$PG_FILE; \
	printf "func New$${PASCAL}Repository(db *gorm.DB) repository.$${PASCAL}Repository {\n" >> $$PG_FILE; \
	printf "\treturn &$${CAMEL}Repository{db: db}\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "func (r *$${CAMEL}Repository) Create(ctx context.Context, v *model.$${PASCAL}) error {\n" >> $$PG_FILE; \
	printf "\treturn r.db.WithContext(ctx).Create(v).Error\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "func (r *$${CAMEL}Repository) GetByID(ctx context.Context, id uint) (*model.$${PASCAL}, error) {\n" >> $$PG_FILE; \
	printf "\tvar v model.$${PASCAL}\n" >> $$PG_FILE; \
	printf "\tif err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {\n" >> $$PG_FILE; \
	printf "\t\treturn nil, err\n" >> $$PG_FILE; \
	printf "\t}\n" >> $$PG_FILE; \
	printf "\treturn &v, nil\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "func (r *$${CAMEL}Repository) List(ctx context.Context, limit, offset int) ([]model.$${PASCAL}, error) {\n" >> $$PG_FILE; \
	printf "\tvar items []model.$${PASCAL}\n" >> $$PG_FILE; \
	printf "\treturn items, r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&items).Error\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "func (r *$${CAMEL}Repository) Update(ctx context.Context, v *model.$${PASCAL}) error {\n" >> $$PG_FILE; \
	printf "\treturn r.db.WithContext(ctx).Save(v).Error\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	printf "\n" >> $$PG_FILE; \
	printf "func (r *$${CAMEL}Repository) Delete(ctx context.Context, id uint) error {\n" >> $$PG_FILE; \
	printf "\treturn r.db.WithContext(ctx).Delete(&model.$${PASCAL}{}, id).Error\n" >> $$PG_FILE; \
	printf "}\n" >> $$PG_FILE; \
	echo "[INFO] Created postgres repo: $$PG_FILE"; \
	PG_REGISTRY=internal/repository/postgres/registry.go; \
	if [ ! -f "$$PG_REGISTRY" ]; then \
		printf "package postgres\n\nimport (\n\t\"$(MODULE)/internal/repository\"\n\n\t\"gorm.io/gorm\"\n)\n\n// Registry holds all postgres repository implementations.\ntype Registry struct {\n\tDB *gorm.DB\n}\n\n// NewRepository creates a new postgres repository registry.\nfunc NewRepository(db *gorm.DB) *Registry {\n\treturn &Registry{DB: db}\n}\n\n// AsRepository converts Registry to repository.Registry.\nfunc (r *Registry) AsRepository() repository.Registry {\n\treturn repository.Registry{\n\t\t// TODO: wire fields here, e.g. CustomerRepo: r.CustomerRepo\n\t}\n}\n" >> $$PG_REGISTRY; \
		echo "[INFO] Created postgres registry: $$PG_REGISTRY"; \
	fi; \
	if grep -q "$${PASCAL}Repo repository.$${PASCAL}Repository" "$$PG_REGISTRY"; then \
		echo "[INFO] $${PASCAL}Repo already in postgres registry"; \
	else \
		awk '/^type Registry struct \{$$/{print;print "\t'"$${PASCAL}"'Repo repository.'"$${PASCAL}"'Repository";next}1' \
			"$$PG_REGISTRY" > "$$PG_REGISTRY.tmp" && mv "$$PG_REGISTRY.tmp" "$$PG_REGISTRY"; \
		echo "[INFO] Added $${PASCAL}Repo to postgres registry"; \
		echo "[INFO] NOTE: init $${PASCAL}Repo: New$${PASCAL}Repository(r.DB) in postgres/registry.go NewRepository()"; \
	fi; \
	echo "[INFO] Repository generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# MIGRATION  (Goose Go-based — stored in internal/migration/)
# Usage : make migration VAL=customers
# Output: internal/migration/0001_customers.go
#         internal/migration/registry.go  (created once as package anchor)
# ─────────────────────────────────────────────────────────────────────────────

migration:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make migration VAL=customers"; \
		exit 1; \
	fi; \
	MIGRATION_DIR=internal/migration; \
	mkdir -p $$MIGRATION_DIR; \
	FILE_COUNT=$$(find $$MIGRATION_DIR -type f -name "*.go" ! -name "registry.go" | wc -l | tr -d ' '); \
	NEXT_NUM=$$((FILE_COUNT + 1)); \
	PADDED=$$(printf "%04d" $$NEXT_NUM); \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	DST=$$MIGRATION_DIR/$${PADDED}_$${LOWER}.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package migration\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"context\"\n" >> $$DST; \
	printf "\t\"database/sql\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"github.com/pressly/goose/v3\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func init() {\n" >> $$DST; \
	printf "\tgoose.AddMigrationContext(up$${PADDED}$${PASCAL}, down$${PADDED}$${PASCAL})\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func up$${PADDED}$${PASCAL}(ctx context.Context, tx *sql.Tx) error {\n" >> $$DST; \
	printf "\t_, err := tx.ExecContext(ctx, \`\n" >> $$DST; \
	printf "\t\tCREATE TABLE IF NOT EXISTS $${LOWER}s (\n" >> $$DST; \
	printf "\t\t\tid         BIGSERIAL    PRIMARY KEY,\n" >> $$DST; \
	printf "\t\t\tcreated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),\n" >> $$DST; \
	printf "\t\t\tupdated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),\n" >> $$DST; \
	printf "\t\t\tdeleted_at TIMESTAMPTZ\n" >> $$DST; \
	printf "\t\t);\n" >> $$DST; \
	printf "\t\`)\n" >> $$DST; \
	printf "\treturn err\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func down$${PADDED}$${PASCAL}(ctx context.Context, tx *sql.Tx) error {\n" >> $$DST; \
	printf "\t_, err := tx.ExecContext(ctx, \`DROP TABLE IF EXISTS $${LOWER}s;\`)\n" >> $$DST; \
	printf "\treturn err\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created migration: $$DST"; \
	REGISTRY_FILE=$$MIGRATION_DIR/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "// Package migration contains Goose Go-based migrations.\n" >> $$REGISTRY_FILE; \
		printf "// Blank-import this package from cmd/server/main.go:\n" >> $$REGISTRY_FILE; \
		printf "//\n" >> $$REGISTRY_FILE; \
		printf "//\timport _ \"$(MODULE)/internal/migration\"\n" >> $$REGISTRY_FILE; \
		printf "package migration\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created migration registry: $$REGISTRY_FILE"; \
	fi; \
	echo "[INFO] Migration generation completed: $${PADDED}_$${LOWER}"

# ─────────────────────────────────────────────────────────────────────────────
# QUEUE PRODUCER
# Usage : make queue-producer VAL=suspend
# Output: internal/queue/producer/suspend_producer.go
#         internal/queue/producer/registry.go  (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

queue-producer:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make queue-producer VAL=suspend"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p internal/queue/producer; \
	DST=internal/queue/producer/$${LOWER}_producer.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package producer\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"context\"\n" >> $$DST; \
	printf "\t\"encoding/json\"\n" >> $$DST; \
	printf "\t\"fmt\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"$(MODULE)/pkg/rabbitmq\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Queue$${PASCAL} is the RabbitMQ queue name for $${LOWER} events.\n" >> $$DST; \
	printf "const Queue$${PASCAL} = \"$${LOWER}\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Event is the message payload published to Queue$${PASCAL}.\n" >> $$DST; \
	printf "type $${PASCAL}Event struct {\n" >> $$DST; \
	printf "\t// TODO: add event fields, e.g. CustomerID uint \`json:\"customer_id\"\`\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Producer publishes $${LOWER} events to RabbitMQ.\n" >> $$DST; \
	printf "type $${PASCAL}Producer struct {\n" >> $$DST; \
	printf "\tclient *rabbitmq.Client\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// New$${PASCAL}Producer creates a new $${LOWER} producer.\n" >> $$DST; \
	printf "func New$${PASCAL}Producer(c *rabbitmq.Client) *$${PASCAL}Producer {\n" >> $$DST; \
	printf "\treturn &$${PASCAL}Producer{client: c}\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Publish publishes a $${LOWER} event.\n" >> $$DST; \
	printf "func (p *$${PASCAL}Producer) Publish(ctx context.Context, exchange, routingKey string, event $${PASCAL}Event) error {\n" >> $$DST; \
	printf "\tbody, err := json.Marshal(event)\n" >> $$DST; \
	printf "\tif err != nil {\n" >> $$DST; \
	printf "\t\treturn fmt.Errorf(\"$${LOWER}_producer: marshal: %%w\", err)\n" >> $$DST; \
	printf "\t}\n" >> $$DST; \
	printf "\treturn p.client.Publish(ctx, exchange, routingKey, body)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created producer: $$DST"; \
	REGISTRY_FILE=internal/queue/producer/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "package producer\n\nimport (\n\t\"$(MODULE)/pkg/rabbitmq\"\n)\n\n// Producer is the registry for all event producers.\ntype Producer struct {\n\tclient *rabbitmq.Client\n}\n\n// NewProducer creates a new producer registry.\nfunc NewProducer(c *rabbitmq.Client) *Producer {\n\treturn &Producer{client: c}\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created producer registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "func (p \*Producer) $${PASCAL}()" "$$REGISTRY_FILE"; then \
		echo "[INFO] $${PASCAL}() already in producer registry"; \
	else \
		printf "\n// $${PASCAL} returns the $${LOWER} producer.\nfunc (p *Producer) $${PASCAL}() *$${PASCAL}Producer {\n\treturn New$${PASCAL}Producer(p.client)\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Added $${PASCAL}() to producer registry"; \
	fi; \
	echo "[INFO] Queue producer generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# QUEUE CONSUMER
# Usage : make queue-consumer VAL=suspend
# Output: internal/queue/consumer/suspend_consumer.go
#         internal/queue/consumer/registry.go  (auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

queue-consumer:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make queue-consumer VAL=suspend"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p internal/queue/consumer; \
	DST=internal/queue/consumer/$${LOWER}_consumer.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package consumer\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"context\"\n" >> $$DST; \
	printf "\t\"encoding/json\"\n" >> $$DST; \
	printf "\t\"fmt\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"$(MODULE)/internal/queue/producer\"\n" >> $$DST; \
	printf "\t\"$(MODULE)/pkg/rabbitmq\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Consumer consumes $${LOWER} events from RabbitMQ.\n" >> $$DST; \
	printf "type $${PASCAL}Consumer struct {\n" >> $$DST; \
	printf "\tclient *rabbitmq.Client\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// New$${PASCAL}Consumer creates a new $${LOWER} consumer.\n" >> $$DST; \
	printf "func New$${PASCAL}Consumer(c *rabbitmq.Client) *$${PASCAL}Consumer {\n" >> $$DST; \
	printf "\treturn &$${PASCAL}Consumer{client: c}\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Start begins consuming messages from the $${LOWER} queue.\n" >> $$DST; \
	printf "// Run this in a goroutine from consumer/registry.go StartAll().\n" >> $$DST; \
	printf "func (c *$${PASCAL}Consumer) Start(ctx context.Context) error {\n" >> $$DST; \
	printf "\treturn c.client.Subscribe(ctx, producer.Queue$${PASCAL}, c.handle)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "func (c *$${PASCAL}Consumer) handle(ctx context.Context, body []byte) error {\n" >> $$DST; \
	printf "\tvar event producer.$${PASCAL}Event\n" >> $$DST; \
	printf "\tif err := json.Unmarshal(body, &event); err != nil {\n" >> $$DST; \
	printf "\t\treturn fmt.Errorf(\"$${LOWER}_consumer: unmarshal: %%w\", err)\n" >> $$DST; \
	printf "\t}\n" >> $$DST; \
	printf "\t// TODO: implement using event data\n" >> $$DST; \
	printf "\t_ = ctx\n" >> $$DST; \
	printf "\treturn nil\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created consumer: $$DST"; \
	REGISTRY_FILE=internal/queue/consumer/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "package consumer\n\nimport (\n\t\"context\"\n\n\t\"$(MODULE)/pkg/rabbitmq\"\n)\n\n// Consumer is the registry for all queue consumers.\ntype Consumer struct {\n\tclient *rabbitmq.Client\n}\n\n// NewConsumer creates a new consumer registry.\nfunc NewConsumer(c *rabbitmq.Client) *Consumer {\n\treturn &Consumer{client: c}\n}\n\n// StartAll launches all consumers in goroutines.\n// Call this once from cmd/server/main.go.\nfunc (c *Consumer) StartAll(ctx context.Context) {\n\t// Register consumers here:\n\t// go func() { _ = c.Suspend().Start(ctx) }()\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created consumer registry: $$REGISTRY_FILE"; \
	fi; \
	if grep -q "func (c \*Consumer) $${PASCAL}()" "$$REGISTRY_FILE"; then \
		echo "[INFO] $${PASCAL}() already in consumer registry"; \
	else \
		printf "\n// $${PASCAL} returns the $${LOWER} consumer.\nfunc (c *Consumer) $${PASCAL}() *$${PASCAL}Consumer {\n\treturn New$${PASCAL}Consumer(c.client)\n}\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Added $${PASCAL}() to consumer registry"; \
	fi; \
	echo "[INFO] NOTE: register goroutine in StartAll → go func() { _ = c.$${PASCAL}().Start(ctx) }()"; \
	echo "[INFO] Queue consumer generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# SCHEDULER
# Usage : make scheduler VAL=billing
# Output: internal/scheduler/billing_scheduler.go
#         internal/scheduler/registry.go  (auto-updated)
# NOTE  : Each scheduler registers cron jobs in its Start() method
# ─────────────────────────────────────────────────────────────────────────────

scheduler:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make scheduler VAL=billing"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p internal/scheduler; \
	DST=internal/scheduler/$${LOWER}_scheduler.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package scheduler\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "import (\n" >> $$DST; \
	printf "\t\"context\"\n" >> $$DST; \
	printf "\t\"log\"\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "\t\"github.com/robfig/cron/v3\"\n" >> $$DST; \
	printf ")\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Scheduler handles $${LOWER} cron jobs.\n" >> $$DST; \
	printf "type $${PASCAL}Scheduler struct {\n" >> $$DST; \
	printf "\tcron *cron.Cron\n" >> $$DST; \
	printf "\t// TODO: inject *service.$${PASCAL}Service and *producer.$${PASCAL}Producer\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// New$${PASCAL}Scheduler creates a new $${LOWER} scheduler.\n" >> $$DST; \
	printf "func New$${PASCAL}Scheduler(c *cron.Cron) *$${PASCAL}Scheduler {\n" >> $$DST; \
	printf "\treturn &$${PASCAL}Scheduler{cron: c}\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// Start registers $${LOWER} cron jobs.\n" >> $$DST; \
	printf "func (s *$${PASCAL}Scheduler) Start() {\n" >> $$DST; \
	printf "\ts.cron.AddFunc(\"@daily\", func() {\n" >> $$DST; \
	printf "\t\tctx := context.Background()\n" >> $$DST; \
	printf "\t\tlog.Println(\"scheduler: running $${LOWER} job\")\n" >> $$DST; \
	printf "\t\t// TODO: implement, e.g. s.service.DoWork(ctx)\n" >> $$DST; \
	printf "\t\t_ = ctx\n" >> $$DST; \
	printf "\t})\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created scheduler: $$DST"; \
	REGISTRY_FILE=internal/scheduler/registry.go; \
	if [ ! -f "$$REGISTRY_FILE" ]; then \
		printf "package scheduler\n\nimport (\n\t\"github.com/robfig/cron/v3\"\n\t\"$(MODULE)/internal/queue\"\n\t\"$(MODULE)/internal/service\"\n)\n\n// Registry holds all scheduler instances.\ntype Registry struct {\n\tcron    *cron.Cron\n\tservice *service.Registry\n\tqueue   *queue.Registry\n}\n\n// NewRegistry creates a new scheduler registry.\nfunc NewRegistry(s *service.Registry, q *queue.Registry) *Registry {\n\tc := cron.New()\n\tr := &Registry{cron: c, service: s, queue: q}\n\t// Register schedulers here:\n\t// NewXxxScheduler(c, s.Xxx, q.XxxProducer).Start()\n\treturn r\n}\n\n// Start starts all cron jobs.\nfunc (r *Registry) Start() { r.cron.Start() }\n\n// Stop stops all cron jobs.\nfunc (r *Registry) Stop() { r.cron.Stop() }\n" >> $$REGISTRY_FILE; \
		echo "[INFO] Created scheduler registry: $$REGISTRY_FILE"; \
	fi; \
	echo "[INFO] NOTE: register in registry.go → New$${PASCAL}Scheduler(r.cron, ...).Start()"; \
	echo "[INFO] Scheduler generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# MIKROTIK DOMAIN ENTITY
# Usage : make mikrotik-domain VAL=ppp
# Output: pkg/mikrotik/domain/ppp.go
# ─────────────────────────────────────────────────────────────────────────────

mikrotik-domain:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make mikrotik-domain VAL=ppp"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
	fi; \
	mkdir -p pkg/mikrotik/domain; \
	DST=pkg/mikrotik/domain/$${LOWER}.go; \
	if [ -f "$$DST" ]; then \
		echo "[ERROR] File $$DST already exists."; \
		exit 1; \
	fi; \
	printf "package mktdomain\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL} represents a Mikrotik $${LOWER} entity.\n" >> $$DST; \
	printf "// Tag \`ros\` maps struct fields to RouterOS API field names.\n" >> $$DST; \
	printf "type $${PASCAL} struct {\n" >> $$DST; \
	printf "\tID       string \`ros:\".id\"\`\n" >> $$DST; \
	printf "\tComment  string \`ros:\"comment\"\`\n" >> $$DST; \
	printf "\tDisabled bool   \`ros:\"disabled\"\`\n" >> $$DST; \
	printf "\t// TODO: add RouterOS-specific fields\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	printf "\n" >> $$DST; \
	printf "// $${PASCAL}Repository defines the contract for RouterOS $${LOWER} operations.\n" >> $$DST; \
	printf "type $${PASCAL}Repository interface {\n" >> $$DST; \
	printf "\tAdd(v $${PASCAL}) error\n" >> $$DST; \
	printf "\tRemove(id string) error\n" >> $$DST; \
	printf "\tEnable(id string) error\n" >> $$DST; \
	printf "\tDisable(id string) error\n" >> $$DST; \
	printf "\tGet(id string) (*$${PASCAL}, error)\n" >> $$DST; \
	printf "\tList() ([]$${PASCAL}, error)\n" >> $$DST; \
	printf "}\n" >> $$DST; \
	echo "[INFO] Created mikrotik domain: $$DST"

# ─────────────────────────────────────────────────────────────────────────────
# MIKROTIK MODULE
# Usage : make mikrotik-module VAL=ppp
# Output: pkg/mikrotik/ppp/service.go
#         pkg/mikrotik/ppp/repository.go
#         pkg/mikrotik/ppp/ppp_test.go
#         pkg/mikrotik/mikrotik.go  (facade — auto-updated)
# ─────────────────────────────────────────────────────────────────────────────

mikrotik-module:
	@if [ -z "$(VAL)" ]; then \
		echo "[ERROR] Please provide VAL, e.g. make mikrotik-module VAL=ppp"; \
		exit 1; \
	fi; \
	LOWER=$$(echo $(VAL) | tr '[:upper:]' '[:lower:]'); \
	if [[ "$$LOWER" == *_* ]]; then \
		PASCAL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)} 1'); \
		CAMEL=$$(echo "$$LOWER" | awk 'BEGIN{FS="_"} {printf "%s", $$1; for(i=2;i<=NF;i++) printf "%s", toupper(substr($$i,1,1)) substr($$i,2)} END{print ""}'); \
	else \
		PASCAL=$$(echo $$LOWER | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}'); \
		CAMEL=$$LOWER; \
	fi; \
	MODULE_DIR=pkg/mikrotik/$$LOWER; \
	if [ -d "$$MODULE_DIR" ]; then \
		echo "[ERROR] Directory $$MODULE_DIR already exists."; \
		exit 1; \
	fi; \
	mkdir -p $$MODULE_DIR; \
	echo "[INFO] Created directory: $$MODULE_DIR"; \
	SVC_FILE=$$MODULE_DIR/service.go; \
	printf "package $${LOWER}\n\nimport mktdomain \"$(MODULE)/pkg/mikrotik/domain\"\n\ntype $${PASCAL}Service interface {\n\tmktdomain.$${PASCAL}Repository\n}\n\ntype $${CAMEL}Service struct {\n\trepo mktdomain.$${PASCAL}Repository\n}\n\nfunc New$${PASCAL}Service(repo mktdomain.$${PASCAL}Repository) $${PASCAL}Service {\n\treturn &$${CAMEL}Service{repo: repo}\n}\n\nfunc (s *$${CAMEL}Service) Add(v mktdomain.$${PASCAL}) error                       { return s.repo.Add(v) }\nfunc (s *$${CAMEL}Service) Remove(id string) error                        { return s.repo.Remove(id) }\nfunc (s *$${CAMEL}Service) Enable(id string) error                        { return s.repo.Enable(id) }\nfunc (s *$${CAMEL}Service) Disable(id string) error                       { return s.repo.Disable(id) }\nfunc (s *$${CAMEL}Service) Get(id string) (*mktdomain.$${PASCAL}, error)  { return s.repo.Get(id) }\nfunc (s *$${CAMEL}Service) List() ([]mktdomain.$${PASCAL}, error)         { return s.repo.List() }\n" >> $$SVC_FILE; \
	echo "[INFO] Created service: $$SVC_FILE"; \
	REPO_FILE=$$MODULE_DIR/repository.go; \
	printf "package $${LOWER}\n\nimport (\n\t\"fmt\"\n\n\t\"$(MODULE)/pkg/mikrotik/client\"\n\tmktdomain \"$(MODULE)/pkg/mikrotik/domain\"\n)\n\ntype $${CAMEL}Repository struct {\n\tclient *client.Client\n}\n\nfunc NewRepository(c *client.Client) mktdomain.$${PASCAL}Repository {\n\treturn &$${CAMEL}Repository{client: c}\n}\n\nfunc (r *$${CAMEL}Repository) Add(v mktdomain.$${PASCAL}) error {\n\t// Example: r.client.Run(\"/$${LOWER}/add\", \"=comment=\"+v.Comment)\n\treturn fmt.Errorf(\"not implemented\")\n}\nfunc (r *$${CAMEL}Repository) Remove(id string) error  { return fmt.Errorf(\"not implemented\") }\nfunc (r *$${CAMEL}Repository) Enable(id string) error  { return fmt.Errorf(\"not implemented\") }\nfunc (r *$${CAMEL}Repository) Disable(id string) error { return fmt.Errorf(\"not implemented\") }\nfunc (r *$${CAMEL}Repository) Get(id string) (*mktdomain.$${PASCAL}, error) {\n\treturn nil, fmt.Errorf(\"not implemented\")\n}\nfunc (r *$${CAMEL}Repository) List() ([]mktdomain.$${PASCAL}, error) {\n\treturn nil, fmt.Errorf(\"not implemented\")\n}\n" >> $$REPO_FILE; \
	echo "[INFO] Created repository: $$REPO_FILE"; \
	TEST_FILE=$$MODULE_DIR/$${LOWER}_test.go; \
	printf "package $${LOWER}_test\n\nimport \"testing\"\n\nfunc Test$${PASCAL}Placeholder(t *testing.T) {\n\tt.Log(\"$${LOWER}: add tests here\")\n}\n" >> $$TEST_FILE; \
	echo "[INFO] Created test file: $$TEST_FILE"; \
	FACADE_FILE=pkg/mikrotik/mikrotik.go; \
	if [ ! -f "$$FACADE_FILE" ]; then \
		mkdir -p pkg/mikrotik; \
		printf "package mikrotik\n\nimport (\n\t\"$(MODULE)/pkg/mikrotik/client\"\n)\n\n// Mikrotik is the single entry point for all RouterOS operations.\n// Add sub-module fields below as you run: make mikrotik-module VAL=<name>\ntype Mikrotik struct {\n}\n\nfunc New(opts client.Options) (*Mikrotik, error) {\n\tc, err := client.New(opts)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\t_ = c\n\treturn &Mikrotik{}, nil\n}\n" >> $$FACADE_FILE; \
		echo "[INFO] Created mikrotik facade: $$FACADE_FILE"; \
	fi; \
	if grep -q "$$PASCAL " "$$FACADE_FILE"; then \
		echo "[INFO] $${PASCAL} already in facade"; \
	else \
		awk '/^type Mikrotik struct \{$$/{print;print "\t'"$${PASCAL} $${LOWER}.$${PASCAL}Service"'";next}1' \
			"$$FACADE_FILE" > "$$FACADE_FILE.tmp" && mv "$$FACADE_FILE.tmp" "$$FACADE_FILE"; \
		echo "[INFO] Added $${PASCAL} field to facade."; \
		echo "[INFO] NOTE: add import for $${LOWER} and init in mikrotik.go manually."; \
	fi; \
	echo "[INFO] Mikrotik module generation completed: $$LOWER"

# ─────────────────────────────────────────────────────────────────────────────
# GENERATE MOCKS
# ─────────────────────────────────────────────────────────────────────────────

generate-mocks:
	@echo "[INFO] Generating mocks from go:generate directives..."
	@mkdir -p tests/mocks/repository tests/mocks/service
	@go generate ./internal/repository/...
	@echo "[INFO] Done."

# ─────────────────────────────────────────────────────────────────────────────
# DEV TOOLS
# ─────────────────────────────────────────────────────────────────────────────

tidy:
	@echo "[INFO] Running go mod tidy..."
	@go mod tidy

lint:
	@echo "[INFO] Running golangci-lint..."
	@golangci-lint run ./...

test:
	@echo "[INFO] Running unit tests..."
	@go test -v -race ./internal/... ./pkg/...

test-coverage:
	@echo "[INFO] Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "[INFO] Coverage report: coverage.html"

test-integration:
	@echo "[INFO] Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...
dev:
	air -c .air.toml

# ─────────────────────────────────────────────────────────────────────────────
# INTERACTIVE TARGET SELECTOR  (fzf recommended: apt/brew install fzf)
# ─────────────────────────────────────────────────────────────────────────────

run:
	@VAL_TARGETS="model domain service handler repository scheduler migration queue-producer queue-consumer mikrotik-domain mikrotik-module"; \
	if command -v fzf >/dev/null 2>&1; then \
		target=$$(grep -E "^[a-zA-Z0-9_-]+:" $(MAKEFILE_LIST) | grep -v "^run:" | sed 's/:.*//' | sort -u | fzf --height=22 --prompt="Select target: "); \
	else \
		echo "[INFO] fzf not found, using basic menu. Install: apt install fzf"; \
		targets=$$(grep -E "^[a-zA-Z0-9_-]+:" $(MAKEFILE_LIST) | grep -v "^run:" | sed 's/:.*//' | sort -u); \
		i=1; for t in $$targets; do echo "$$i) $$t"; i=$$((i+1)); done; \
		printf "Enter number: "; read -r choice; \
		target=$$(echo "$$targets" | sed -n "$${choice}p"); \
	fi; \
	[ -z "$$target" ] && echo "[INFO] No target selected." && exit 0; \
	echo "[INFO] Selected: $$target"; \
	needs_val=0; \
	for vt in $$VAL_TARGETS; do [ "$$target" = "$$vt" ] && needs_val=1 && break; done; \
	if [ "$$needs_val" = "1" ]; then \
		printf "Enter VAL: "; val=$$(bash -c 'read -r v && echo "$$v"'); \
		[ -z "$$val" ] && echo "[ERROR] VAL required for: $$target" && exit 1; \
		$(MAKE) $$target VAL=$$val; \
	elif [ "$$target" = "http" ]; then \
		printf "Force rebuild? (y/N): "; build=$$(bash -c 'read -r b && echo "$$b"'); \
		if [ "$$build" = "y" ] || [ "$$build" = "Y" ]; then $(MAKE) http BUILD=true; else $(MAKE) http; fi; \
	else \
		$(MAKE) $$target; \
	fi

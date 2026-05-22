APP_NAME   := library-cli
ENTRYPOINT := ./cmd/app
BUILD_DIR  := ./bin

.PHONY: run build test clean tidy fmt vet lint help
.PHONY: db-create db-migrate db-seed db-reset

## ── Development ─────────────────────────────────────────────

run: build
	@echo "🚀 Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

build: tidy
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(ENTRYPOINT)

test:
	@echo "🧪 Running tests..."
	@go test ./tests/ -v -count=1

test-cover:
	@echo "🧪 Running tests with coverage..."
	@go test ./tests/ -coverprofile=coverage.out && go tool cover -html=coverage.out

## ── Quality ──────────────────────────────────────────────────

tidy:
	@echo "📦 Tidying dependencies..."
	@go mod tidy

fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

## ── Database ─────────────────────────────────────────────────

db-create:
	@echo "🗄️  Creating database..."
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS library_db;"

db-migrate:
	@echo "📜 Running schema migration..."
	mysql -u root -p library_db < schema.sql

db-seed:
	@echo "🌱 Seeding database..."
	mysql -u root -p library_db < seed.sql

db-reset: db-create db-migrate db-seed
	@echo "✅ Database reset complete."

## ── Cleanup ──────────────────────────────────────────────────

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) coverage.out

## ── Help ─────────────────────────────────────────────────────

help:
	@printf "┌──────────────────────────────────────────┐\n"
	@printf "│  📚 $(APP_NAME) — Makefile Commands      │\n"
	@printf "├──────────────────────────────────────────┤\n"
	@printf "│  make run          Build & run app       │\n"
	@printf "│  make build        Compile binary        │\n"
	@printf "│  make test         Run tests             │\n"
	@printf "│  make test-cover   Tests + coverage HTML │\n"
	@printf "│  make tidy         go mod tidy           │\n"
	@printf "│  make fmt          Format code           │\n"
	@printf "│  make vet          Static analysis       │\n"
	@printf "│  make lint         Run golangci-lint     │\n"
	@printf "│  make db-create    Create MySQL DB       │\n"
	@printf "│  make db-migrate   Run schema.sql        │\n"
	@printf "│  make db-seed      Run seed.sql          │\n"
	@printf "│  make db-reset     Reset DB              │\n"
	@printf "│  make clean        Remove binaries       │\n"
	@printf "└──────────────────────────────────────────┘\n"

.DEFAULT_GOAL := help

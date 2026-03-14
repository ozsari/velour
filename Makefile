.PHONY: build dev frontend backend docker clean

# Build everything
build: frontend backend

# Development mode
dev:
	@echo "Starting backend..."
	@go run ./cmd/velour &
	@echo "Starting frontend dev server..."
	@cd web/frontend && npm run dev

# Build frontend
frontend:
	cd web/frontend && npm ci --legacy-peer-deps && npm run build

# Build Go backend
backend:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o velour ./cmd/velour

# Docker build
docker:
	docker compose build

# Docker run
docker-up:
	docker compose up -d

# Docker stop
docker-down:
	docker compose down

# Clean
clean:
	rm -f velour
	rm -rf web/dist
	rm -rf web/frontend/node_modules

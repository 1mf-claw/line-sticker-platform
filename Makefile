.PHONY: lint test build

lint:
	cd frontend && npm run lint || true

test:
	cd backend && go test ./...

build-backend:
	mkdir -p dist/backend
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../dist/backend/app-linux-amd64 ./cmd/app

build-frontend:
	cd frontend && npm ci && npm run build
	mkdir -p dist/frontend
	cp -r frontend/dist/* dist/frontend/

build: build-backend build-frontend

fmt:
	go fmt ./...

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

run-gateway:
	go run ./cmd/gateway

run-backend:
	go run ./cmd/mock

run-auth:
	go run ./cmd/auth

run-auth-postgres:
	docker compose -f deployments/docker-compose.yml up -d
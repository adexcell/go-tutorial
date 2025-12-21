.PHONY:
start:
	docker-compose up --build && curl http://localhost:8080/health

.PHONY:
local:
	go run cmd/server/main.go
.PHONY: build
build: ## Compilation
	go build -o ./bin/evf ./cmd/evf/main.go

.PHONY: run
run: ## Runs the binary
	go run cmd/evf/main.go
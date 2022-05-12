.PHONY: build
build: ## Compilation
	go build -o ./bin/evf ./cmd/evf/main.go

.PHONY: run
run: ## Runs the binary
	go run cmd/evf/main.go

# Run the unit tests
.PHONY: test
test: ## Run the unit tests
	go test -v -coverprofile cover.out ./...
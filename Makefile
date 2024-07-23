.PHONY: init
init:
	@goimports -local "github.com/chiuchungho/eth-usdc-transfer-exporter" -w .
	@go mod vendor

.PHONY: compose_up ## Downloads dependencies, runs tests, builds image based with test tag, and runs docker-compose
compose_up: test image_test
	@docker-compose up

.PHONY: image_test
image_test: build_image_test

# Testing targets
.PHONY: build_image_test
build_image_test:
	@echo "Building Docker Image test"
	@docker build . -t eth-usdc-transfer-exporter:test --file deploy/eth-usdc-transfer-exporter.dockerfile --no-cache

.PHONY: test
test: ## Run all tests
	$(call blue, "# Running tests...")
	go test ./... -tags=integration

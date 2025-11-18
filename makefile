check: ## Run go vet and go fmt
	@echo "Running go fmt..."
	@unformatted=$$(go fmt ./...); \
	if [ -n "$$unformatted" ]; then \
		echo "❌ Unformatted files found:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ Formatting and vetting passed"

migrationup:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/iot-asset-tracking?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/iot-asset-tracking?sslmode=disable" -verbose down

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total

coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-check:
	go test -coverprofile=coverage.out ./internal/fetcher ./internal/fetcher/strategies ./internal/model > /dev/null
	go tool cover -func=coverage.out | grep github.com/wosiu6/docky-go/internal/fetcher | awk '{print $$3}' | sed 's/%//' | awk '{ if ($$1 < 80.0) { printf "Fetcher coverage %.1f%% below threshold\n", $$1; exit 1 } }'
	@echo "Coverage thresholds met"
go test -coverprofile=coverage.out ./...
if (Test-Path coverage.out) {
  go tool cover -func=coverage.out | Select-String 'total'
}

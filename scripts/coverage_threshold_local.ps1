param(
  [double]$Min = 80.0
)
if (-not (Test-Path coverage.out)) { Write-Host 'coverage.out not found. Run go test -coverprofile=coverage.out ./... first.'; exit 2 }
$lines = go tool cover -func=coverage.out 2>$null
$fetcher = $lines -split "`n" | Where-Object { $_ -match 'github.com/wosiu6/docky-go/internal/fetcher$' }
if (-not $fetcher) { Write-Host 'Fetcher line not found'; exit 3 }
$percent = [double]($fetcher.Trim().Split()[-1].TrimEnd('%'))
if ($percent -lt $Min) { Write-Host "Fetcher coverage $percent% below threshold $Min%"; exit 1 }
Write-Host "Fetcher coverage $percent% OK (>= $Min%)"
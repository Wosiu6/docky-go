$minFetcher = 80.0
$pkgs = @('internal/fetcher','internal/fetcher/strategies','internal/model')
$null = go test -coverprofile=coverage.out $pkgs 2>$null
$func = go tool cover -func=coverage.out
$fetcherLine = ($func -split "`n" | Select-String 'internal/fetcher$')
if (-not $fetcherLine) { Write-Host 'Fetcher package not found in coverage output'; exit 2 }
$percent = [double]($fetcherLine.ToString().Split()[-1].TrimEnd('%'))
if ($percent -lt $minFetcher) { Write-Host "Fetcher coverage $percent% below threshold $minFetcher%"; exit 3 }
Write-Host "Fetcher coverage $percent% >= $minFetcher% : OK"
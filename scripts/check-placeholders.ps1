$found = Select-String `
  -Path ".\**\*.json",".\**\*.yml",".\**\*.yaml",".\**\*.md",".\**\*.env*",".\Makefile",".\Dockerfile" `
  -Pattern "__[A-Z_]+__" -Recurse |
  Where-Object {
    $_.Path -notlike "*template.config.json*" -and
    $_.Path -notlike "*check-placeholders*"
  }

if ($found.Count -gt 0) {
  Write-Host "FAIL: $($found.Count) placeholder(s) remaining:" -ForegroundColor Red
  $found | ForEach-Object { Write-Host "  $($_.Path):$($_.LineNumber)  $($_.Line.Trim())" }
  exit 1
} else {
  Write-Host "OK: no placeholders remaining" -ForegroundColor Green
}

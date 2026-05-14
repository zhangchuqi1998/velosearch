# Day 11 crash recovery test.
#
# Usage:
#   .\benchmark\crash_test.ps1                # 10 iterations, n=1000
#   .\benchmark\crash_test.ps1 -Iterations 3  # quick smoke
#
# Each iteration:
#   1. Start server with a fresh data-dir
#   2. Client writes N deterministic vectors (mode=write)
#   3. Force-kill the server (Windows analog of SIGKILL)
#   4. Restart server pointing at the same data-dir (triggers WAL replay)
#   5. Client searches for each id (mode=verify); every id must be its own top-1
#
# Failed iterations leave their data-dir + server logs behind for inspection.

param(
    [int]$Iterations = 10,
    [int]$N = 1000,
    [string]$Addr = ":50052",
    [string]$ClientAddr = "localhost:50052"
)

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$binDir = Join-Path $repoRoot "bin"
$serverExe = Join-Path $binDir "velosearch.exe"
$clientExe = Join-Path $binDir "crash_client.exe"

# --- Build once -------------------------------------------------------------

Write-Host "Building binaries..."
Push-Location $repoRoot
try {
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir | Out-Null
    }
    & go build -o $serverExe ./cmd/server
    if ($LASTEXITCODE -ne 0) { throw "build server failed (exit $LASTEXITCODE)" }
    & go build -o $clientExe ./benchmark/crash_client
    if ($LASTEXITCODE -ne 0) { throw "build crash_client failed (exit $LASTEXITCODE)" }
}
finally {
    Pop-Location
}

# --- Iteration --------------------------------------------------------------

function Invoke-CrashIteration {
    param([int]$Iter)

    $dataDir = Join-Path $repoRoot "testdata\crash_$([guid]::NewGuid().Guid)"
    New-Item -ItemType Directory -Force -Path $dataDir | Out-Null

    $success = $false
    $serverProc = $null
    $server2Proc = $null

    try {
        Write-Host "[Run $Iter, 1/5] Start server..."
        $serverProc = Start-Process -FilePath $serverExe `
            -ArgumentList "-data-dir=$dataDir","-addr=$Addr" `
            -PassThru `
            -RedirectStandardOutput "$dataDir\server.log" `
            -RedirectStandardError  "$dataDir\server.err"
        Start-Sleep -Seconds 2
        if ($serverProc.HasExited) {
            throw "server died on startup; see $dataDir\server.err"
        }

        Write-Host "[Run $Iter, 2/5] Write $N vectors..."
        & $clientExe "-addr=$ClientAddr" "-mode=write" "-n=$N"
        if ($LASTEXITCODE -ne 0) { throw "write failed (exit $LASTEXITCODE)" }

        Write-Host "[Run $Iter, 3/5] SIGKILL server..."
        Stop-Process -Id $serverProc.Id -Force
        $serverProc.WaitForExit(5000) | Out-Null
        Start-Sleep -Seconds 1

        Write-Host "[Run $Iter, 4/5] Restart server (triggers replay)..."
        $server2Proc = Start-Process -FilePath $serverExe `
            -ArgumentList "-data-dir=$dataDir","-addr=$Addr" `
            -PassThru `
            -RedirectStandardOutput "$dataDir\server2.log" `
            -RedirectStandardError  "$dataDir\server2.err"
        Start-Sleep -Seconds 2
        if ($server2Proc.HasExited) {
            throw "server died on restart; see $dataDir\server2.err"
        }

        Write-Host "[Run $Iter, 5/5] Verify all $N ids..."
        & $clientExe "-addr=$ClientAddr" "-mode=verify" "-n=$N"
        $verifyExit = $LASTEXITCODE

        Stop-Process -Id $server2Proc.Id -Force
        $server2Proc.WaitForExit(5000) | Out-Null

        if ($verifyExit -ne 0) {
            throw "verify failed (exit $verifyExit); see $dataDir\server2.log"
        }

        Write-Host "[Run $Iter] OK" -ForegroundColor Green
        $success = $true
    }
    catch {
        Write-Warning "[Run $Iter] FAILED: $_"
        Write-Warning "Data preserved: $dataDir"
        # Make sure no orphan processes are left behind
        foreach ($p in @($serverProc, $server2Proc)) {
            if ($p -and -not $p.HasExited) {
                try { Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue } catch {}
            }
        }
    }

    if ($success) {
        Remove-Item -Recurse -Force $dataDir -ErrorAction SilentlyContinue
    }
    return $success
}

# --- Loop -------------------------------------------------------------------

$passed = 0
for ($i = 1; $i -le $Iterations; $i++) {
    if (Invoke-CrashIteration -Iter $i) { $passed++ }
}

Write-Host ""
Write-Host "================================================"
Write-Host "Result: $passed / $Iterations passed"
Write-Host "================================================"

if ($passed -ne $Iterations) { exit 1 }
exit 0

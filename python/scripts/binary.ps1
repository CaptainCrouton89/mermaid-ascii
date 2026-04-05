$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$PythonDir = Split-Path -Parent $ScriptDir
$GoRoot = Split-Path -Parent $PythonDir
$Target = Join-Path $PythonDir "src\mermaid_ascii\mermaid-ascii.exe"

Push-Location $GoRoot
try {
    go build -o $Target .
} finally {
    Pop-Location
}

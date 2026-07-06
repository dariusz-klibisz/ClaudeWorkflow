# wf bootstrap (SessionStart, native Windows): install the platform engine
# binary into ${CLAUDE_PLUGIN_DATA}\bin\wf.exe — the stable hook path (07 §4).
$ErrorActionPreference = "SilentlyContinue"

$root = $env:CLAUDE_PLUGIN_ROOT
$data = $env:CLAUDE_PLUGIN_DATA
if (-not $root -or -not $data) { exit 0 }

$want = ""
if (Test-Path "$root\bin\VERSION") { $want = (Get-Content "$root\bin\VERSION" -Raw).Trim() }
$have = ""
if (Test-Path "$data\bin\VERSION") { $have = (Get-Content "$data\bin\VERSION" -Raw).Trim() }
if ($want -and ($want -eq $have) -and (Test-Path "$data\bin\wf.exe")) { exit 0 }

$arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $arch = "arm64" }
$src = "$root\bin\wf-windows-$arch.exe"
if (-not (Test-Path $src)) {
  Write-Error "[wf bootstrap] no engine binary for windows/$arch under $root\bin — wf gates will fail open"
  exit 0
}

# checksum verification when the sums file ships
$sums = "$root\bin\SHA256SUMS"
if (Test-Path $sums) {
  $name = Split-Path $src -Leaf
  $line = Select-String -Path $sums -Pattern ([regex]::Escape($name)) | Select-Object -First 1
  if ($line) {
    $expected = ($line.Line -split '\s+')[0].ToLower()
    $actual = (Get-FileHash -Algorithm SHA256 -Path $src).Hash.ToLower()
    if ($expected -ne $actual) {
      Write-Error "[wf bootstrap] checksum mismatch for $name — refusing to install"
      exit 0
    }
  }
}

New-Item -ItemType Directory -Force -Path "$data\bin" | Out-Null
Copy-Item $src "$data\bin\wf.exe" -Force
Copy-Item $src "$data\bin\wf" -Force
if ($want) { Set-Content -Path "$data\bin\VERSION" -Value $want -NoNewline }
Write-Error "[wf bootstrap] installed wf (windows/$arch) -> $data\bin\wf.exe"
exit 0

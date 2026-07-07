# wf bootstrap (SessionStart, native Windows): install the platform engine
# binary into ${CLAUDE_PLUGIN_DATA}\bin\wf.exe — the stable hook path (07 §4).
$ErrorActionPreference = "SilentlyContinue"

$root = $env:CLAUDE_PLUGIN_ROOT
$data = $env:CLAUDE_PLUGIN_DATA
if (-not $root -or -not $data) { exit 0 }

$arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $arch = "arm64" }
$src = "$root\bin\wf-windows-$arch.exe"
$name = "wf-windows-$arch.exe"

# Fetch tier (07 §4-B): no bundled binary, but the committed bin/MANIFEST
# names a released one — download and verify against the COMMITTED checksum.
$manifest = "$root\bin\MANIFEST"
if (-not (Test-Path $src) -and (Test-Path $manifest)) {
  $lines = Get-Content $manifest
  $mver = ($lines | Where-Object { $_ -match '^version ' } | Select-Object -First 1) -replace '^version ', ''
  $murl = ($lines | Where-Object { $_ -match '^base_url ' } | Select-Object -First 1) -replace '^base_url ', ''
  $mline = $lines | Where-Object { $_ -match [regex]::Escape(" $name") + '$' } | Select-Object -First 1
  $msum = if ($mline) { ($mline -split '\s+')[0].ToLower() } else { "" }
  if ($mver -and $murl -and $msum) {
    # exactly this release already installed? — no refetch when bin/ is read-only
    $short = "sha256:" + $msum.Substring(0, 16)
    $have = ""
    if (Test-Path "$data\bin\VERSION") { $have = (Get-Content "$data\bin\VERSION" -Raw).Trim() }
    if (($short -eq $have) -and (Test-Path "$data\bin\wf.exe")) { exit 0 }
    $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("wf-fetch-" + [System.IO.Path]::GetRandomFileName())
    try { Invoke-WebRequest -Uri "$murl/v$mver/$name" -OutFile $tmp -UseBasicParsing } catch {}
    if (Test-Path $tmp) {
      $got = (Get-FileHash -Algorithm SHA256 -Path $tmp).Hash.ToLower()
      if ($got -eq $msum) {
        try { Copy-Item $tmp $src -Force; Write-Error "[wf bootstrap] fetched $name v$mver (checksum verified)" }
        catch { $src = $tmp; $tmp = ""; Write-Error "[wf bootstrap] fetched $name v$mver (checksum verified; plugin bin\ read-only)" }
      } else {
        Write-Error "[wf bootstrap] checksum mismatch for fetched $name (expected $msum, got $got) — refusing"
      }
      if ($tmp -and (Test-Path $tmp)) { Remove-Item $tmp -Force }
    } else {
      Write-Error "[wf bootstrap] could not fetch $murl/v$mver/$name — falling back"
    }
  }
}

if (-not (Test-Path $src)) {
  Write-Error "[wf bootstrap] no engine binary for windows/$arch (bundled or fetched) — wf gates will fail open"
  exit 0
}

# Version stamp: bin/VERSION when it ships; otherwise fall back to the
# binary's checksum so git installs (which .gitignore the VERSION file)
# still get no-op re-runs and a written $data\bin\VERSION.
$want = ""
if (Test-Path "$root\bin\VERSION") { $want = (Get-Content "$root\bin\VERSION" -Raw).Trim() }
if (-not $want) {
  $want = "sha256:" + (Get-FileHash -Algorithm SHA256 -Path $src).Hash.ToLower().Substring(0, 16)
}
$have = ""
if (Test-Path "$data\bin\VERSION") { $have = (Get-Content "$data\bin\VERSION" -Raw).Trim() }
if ($want -and ($want -eq $have) -and (Test-Path "$data\bin\wf.exe")) { exit 0 }

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
# A RUNNING wf.exe cannot be overwritten (the engine self-update path
# replaces the binary that invoked this script) — rename it aside first;
# the .old copies get reaped on the next run.
Remove-Item "$data\bin\wf.exe.old", "$data\bin\wf.old" -Force -ErrorAction SilentlyContinue
if (Test-Path "$data\bin\wf.exe") { Move-Item "$data\bin\wf.exe" "$data\bin\wf.exe.old" -Force }
if (Test-Path "$data\bin\wf") { Move-Item "$data\bin\wf" "$data\bin\wf.old" -Force }
Copy-Item $src "$data\bin\wf.exe" -Force
Copy-Item $src "$data\bin\wf" -Force
if ($want) { Set-Content -Path "$data\bin\VERSION" -Value $want -NoNewline }
Write-Error "[wf bootstrap] installed wf (windows/$arch) -> $data\bin\wf.exe"
exit 0

# remove the space in the following to require running this script as administrator
# Requires -RunAsAdministrator

# terminate on uncaught exceptions
$ErrorActionPreference = "Stop"

# Configure supported HTTPS protocols
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12, [Net.SecurityProtocolType]::Tls11, [Net.SecurityProtocolType]::Tls

$libVersion = '0.15.1'
$libVariant = 'objectbox'  # or 'objectbox-sync'
$downloadDir = 'download'
$extractedLibDir = "$downloadDir\objectbox-$libVersion"

function DownloadLibrary {
    # Windows hashes and URL sourced from objectbox-c download.sh
    $machine = if ([Environment]::Is64BitOperatingSystem) { 'x64' } else { 'x86' }
    $url = "https://github.com/objectbox/objectbox-c/releases/download/v$libVersion/$libVariant-windows-$machine.zip"

    $archiveFile = "$downloadDir\objectbox-$libVersion.zip"

    if (!(Test-Path "$downloadDir" -PathType Container)) {
        Write-Host "Creating download directory: '$downloadDir'"
        New-Item "$downloadDir" -ItemType directory
    }

    Write-Host "Downloading C-API v$libVersion into $archiveFile"
    Write-Host "Downloading from URL: $url"

    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile($url, $archiveFile)

    if (!(Test-Path "$extractedLibDir" -PathType Container)) {
        New-Item "$extractedLibDir" -ItemType directory
    }

    Write-Host "Extracting into $extractedLibDir"
    Expand-Archive "$archiveFile" "$extractedLibDir" -Force
}

function ValidateInstallation {
    param ($sourceFile, $targetFile)

    if ((Get-Item $sourceFile).Length -eq (Get-Item $targetFile).Length) {
        Write-Host "Successfully installed $targetFile" -ForegroundColor Green
    } else {
        throw "Installation to $targetFile failed - source and target contents don't match"
    }
}

function InstallWithPrompt {
    param ($libDir)

    $reply = Read-Host -Prompt "Would you like to install ObjectBox library into $($libDir)? [y/N]"
    if (!($reply -match "[yY]")) {
        Write-Host "OK, skipping installation to $libDir"
        return
    }

    $sourceFile = "$extractedLibDir\lib\objectbox.dll"
    $targetFile = "$libDir\objectbox.dll"

    Write-Host "Copying $sourceFile to $libDir"
    try {
        Copy-Item $sourceFile $libDir
        ValidateInstallation $sourceFile $targetFile
        return
    } catch [System.UnauthorizedAccessException] {
        Write-Host "Can't copy: $($_.Exception.Message)" -ForegroundColor Yellow
    }

    # reaches here only when copying fails because of UnauthorizedAccessException
    $reply = Read-Host -Prompt "Would you like to retry as administrator? [y/N]"
    if (!($reply -match "[yY]")) {
        Write-Host "OK, skipping installation to $libDir"
        return
    }

    $sourceFile = "$pwd\$sourceFile" # sub-shell requires an absolute path. -WorkingDirectory argument doesn't work either.
    $expectedSize = (Get-Item $sourceFile).Length
    $verifyCmd = "if ((Get-Item $targetFile).Length -ne $expectedSize) {Write-Host 'Installation failed.'; Read-Host -Prompt 'Press any key to exit this window'}"
    $cmd = "Copy-Item $sourceFile $libDir ; $verifyCmd"
    Start-Process powershell.exe -Verb runas -ArgumentList $cmd

    ValidateInstallation $sourceFile $targetFile
}

function InstallIntoGCC {
    Write-Host "Determining path to your local GCC installation - necessary for compiling your programs with ObjectBox"

    $libDir = "."
    try {
        # try to find gcc
        $gcc = Get-Command gcc
        Write-Host "Found GCC: $($gcc.Path)"
        $libDir = Split-Path $gcc.Path -Parent
        $libDir = Split-Path $libDir -Parent
        $libDir = Join-Path $libDir "lib"
    }
    catch {
        Write-Host "GCC installation not found, skipping"
        return
    }

    InstallWithPrompt $libDir
}

function InstallIntoSys32 {
    Write-Host "Windows needs to find your library during runtime (incl. tests execution)."
    Write-Host "The simplest way to achieve this is to install the library globally."
    InstallWithPrompt "C:\Windows\System32"
}

try {
    DownloadLibrary
    Write-Host ""

    InstallIntoGCC
    Write-Host ""

    InstallIntoSys32
    Write-Host ""

    Write-Host "Installation complete" -ForegroundColor Green
} catch {
    Write-Error $error[0].Exception -ErrorAction Continue
}

Read-Host -Prompt 'Press any key to exit'
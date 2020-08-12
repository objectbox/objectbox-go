# remove the space in the following to require running this script as administrator
# Requires -RunAsAdministrator

# terminate on uncaught exceptions
$ErrorActionPreference = "Stop"

# Configure supported HTTPS protocols
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12, [Net.SecurityProtocolType]::Tls11, [Net.SecurityProtocolType]::Tls

function DownloadLibrary {
    $libVersion = '0.10.0'
    $downloadDir = 'download'
    
    # Windows hashes and URL sourced from objectbox-c download.sh
    $hash = if ([Environment]::Is64BitOperatingSystem) { "ca33edce272a279b24f87dc0d4cf5bbdcffbc187" } else { "11e6a84a7894f41df553e7c92534c3bf26896802" }
    $remoteRepo = 'https://dl.bintray.com/objectbox/conan/objectbox/objectbox-c'
    $repoType = 'testing'
    $url = "$remoteRepo/$libVersion/$repoType/0/package/$hash/0/conan_package.tgz"
    
    $archiveFile = "$downloadDir\objectbox-$libVersion.tgz"
    $targetDir = "$downloadDir\objectbox-$libVersion"

    if (!(Test-Path $downloadDir -PathType Container)) {
        Write-Host "Creating download directory $downloadDir"
        New-Item $downloadDir -ItemType directory
    }

    Write-Host "Downloading C-API v$libVersion into $archiveFile"

    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile($url, $archiveFile)

    if (!(Test-Path $targetDir -PathType Container)) {
        New-Item $targetDir -ItemType directory 
    }
    
    Write-Host "Extracting into $targetDir"
    tar -xzf "$archiveFile" -C "$targetDir"

    # ls -r
    # Get-ChildItem -Recurse $targetDir 

    return $targetDir
}

function ValdiateInstallation {
    param ($sourceFile, $targetFile)
    
    if ((Get-Item $sourceFile).Length -eq (Get-Item $targetFile).Length) {
        Write-Host "Succesfully installed $targetFile" -ForegroundColor Green
    } else {
        throw "Installation to $targetFile failed - source and target contents don't match"
    }
}

function InstallWithPrompt {
    param ($downloadDir, $libDir)

    $reply = Read-Host -Prompt "Would you like to install ObjectBox library into $($libDir)? [y/N]"  
    if (!($reply -match "[yY]")) { 
        Write-Host "OK, skipping installation to $libDir"
        return
    }
    
    $sourceFile = "$downloadDir\lib\objectbox.dll"
    $targetFile = "$libDir\objectbox.dll"
    
    Write-Host "Copying $sourceFile to $libDir"
    try {
        Copy-Item $sourceFile $libDir
        ValdiateInstallation $sourceFile $targetFile
        return
    } catch [System.UnauthorizedAccessException] {
        Write-Host "Can't copy: $($_.Exception.Message)" -ForegroundColor Yellow
    }

    # reaches here only when copying fails because of UnauthorizedAccessException
    $reply = Read-Host -Prompt "Would you like to retry as adminstrator? [y/N]" 
    if (!($reply -match "[yY]")) {
        Write-Host "OK, skipping installation to $libDir"
        return
    }
    
    $sourceFile = "$pwd\$sourceFile" # sub-shell requires an asbolute path. -WorkingDirectory argument doesn't work either.
    $expectedSize = (Get-Item $sourceFile).Length
    $verifyCmd = "if ((Get-Item $targetFile).Length -ne $expectedSize) {Write-Host 'Installation failed.'; Read-Host -Prompt 'Press any key to exit this window'}"
    $cmd = "Copy-Item $sourceFile $libDir ; $verifyCmd"
    Start-Process powershell.exe -Verb runas -ArgumentList $cmd 

    ValdiateInstallation $sourceFile $targetFile
}

function InstallIntoGCC {
    param ($downloadDir)

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

    InstallWithPrompt $downloadDir $libDir
}

function InstallIntoSys32 {
    param ($downloadDir)

    Write-Host "Windows needs to find your library during runtime (incl. tests execution)." 
    Write-Host "The simplest way to achieve this is to install the library globally."
    InstallWithPrompt $downloadDir "C:\Windows\System32"
}

try {
    $downloadDir = DownloadLibrary
    Write-Host ""

    InstallIntoGCC $downloadDir
    Write-Host ""

    InstallIntoSys32 $downloadDir
    Write-Host ""

    Write-Host "Installation complete" -ForegroundColor Green
} catch {
    Write-Error $error[0].Exception -ErrorAction Continue
}

Read-Host -Prompt 'Press any key to exit'
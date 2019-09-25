# terminate on exceptions
$ErrorActionPreference = "Stop"

# Configure supported protocols
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12, [Net.SecurityProtocolType]::Tls11, [Net.SecurityProtocolType]::Tls

#Write-Host "Following relative paths start within the current directory $pwd"

function DownloadLibrary {
    $libVersion = '0.6.0'
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

function InstallWithPrompt {
    param ($downloadDir, $libDir)

    $reply = Read-Host -Prompt "Would you like to install ObjectBox library into $($libDir)? [y/n]"  
    if (!($reply -match "[yY]")) { 
        Write-Host "OK, skipping installation to $libDir"
        return
    }
    
    $srcFile = "$downloadDir\lib\objectbox.dll"
    Write-Host "Copying $srcFile to $libDir"
    Copy-Item $srcFile $libDir
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

    Write-Host "Windows needs to find your library during runtime (incl. tests execution). The simplest way to achieve this is to install the library globally."
    InstallWithPrompt $downloadDir "C:\Windows\System32"
}

$downloadDir = DownloadLibrary
InstallIntoGCC $downloadDir
InstallIntoSys32 $downloadDir

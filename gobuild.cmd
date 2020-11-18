@echo off

set GOARCH=amd64
set GOOS=windows
set RELEASE=0
set GENERATE=0
set LINUX=0
set MODULE=awsman.exe

if /I "%1"=="release" (
    set RELEASE=1
    echo build release.
)

if /I "%1"=="linux" (
    echo build linux.
    set LINUX=1
    set MODULE=awsman
    
    if /I "%2"=="release" (
        set RELEASE=1
        echo build release.
    )
)

if not exist build.json (
    echo not found build.json
    goto :EOF
)

echo clean ...
go clean
::if exist resource.syso (
::del resource.syso 2>&1
::)

if not exist go.mod (
    echo golang mod init...
    go mod init
)

::run go generate?
if not exist resource.syso SET GENERATE=1
if not exist version.go SET GENERATE=1
if %RELEASE% equ 1 SET GENERATE=1

::go generate
if %GENERATE% equ 1 (
    :: go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
    echo generate resource...
    go generate
    if %errorlevel% neq 0 (
        echo failed.
        goto :EOF
    )
)

::change file time
set filetime=
if exist created.txt (
set /p filetime=< created.txt
)

::echo test ...
::go test

echo build ...
::go build -ldflags "-s -w -extldflags -static -X 'main.version=%version%'" -a -i .
::go build -ldflags "-s -w -linkmode=internal" -a -i -o %MODULE% .
if %LINUX% equ 1 (
    ::linux
    set GOARCH=amd64
    set GOOS=linux
    go build -ldflags "-s -w" -a -o %MODULE% .
) else (
    ::windows
    go build -ldflags "-s -w" -o %MODULE% .
)

if %errorlevel% equ 0 (
    if "%filetime%"=="" (
        echo done.
    ) else (
        if %RELEASE% equ 1 (
            echo compress ...
            upx --version > nul 2>&1
            if %errorlevel% equ 0 (
                upx -9 %MODULE% > nul 2>&1
                if %errorlevel% equ 0 (
                    echo failed to compress module.
                )
            )

            echo change built timestamp ...
            ::echo set LastWriteTime to %filetime% ...
            powershell Set-ItemProperty "%~dp0\%MODULE%" -Name LastWriteTime -Value '%filetime%'
            ::echo set CreationTime  to %filetime% ...
            powershell Set-ItemProperty "%~dp0\%MODULE%" -Name CreationTime  -Value '%filetime%' 
        )
        echo done.
    )
) else (
    echo failed.
    goto :EOF
)

@echo off
setlocal enabledelayedexpansion

REM Get the directory to scan from parameter or use current directory
set "SCAN_DIR=%~1"
if "%SCAN_DIR%"=="" set "SCAN_DIR=."

REM Convert to absolute path
pushd "%SCAN_DIR%" 2>nul
if errorlevel 1 (
    echo Error: Directory "%SCAN_DIR%" does not exist
    exit /b 1
)
set "SCAN_DIR=%CD%"
popd

echo Finding and fixing unnecessary type arguments...
echo Scanning directory: %SCAN_DIR%
echo.

REM Collect unique package directories (one gopls check call per dir, not per file)
set "PREV_DIR="
for /r "%SCAN_DIR%" %%f in (*.go) do (
    set "CURR_DIR=%%~dpf"
    set "CURR_DIR=!CURR_DIR:~0,-1!"
    if not "!CURR_DIR!"=="!PREV_DIR!" (
        set "PREV_DIR=!CURR_DIR!"
        echo Checking !CURR_DIR!...

        REM Run gopls check on the whole directory at once
        for /f "tokens=1,2,3,4 delims=:" %%a in ('gopls check --severity=hint "!CURR_DIR!" 2^>^&1 ^| findstr /C:"unnecessary type arguments" ^| sort /R') do (
            set "drive=%%a"
            set "file=%%b"
            set "line=%%c"
            set "col=%%d"

            REM Construct position
            set "pos=!drive!:!file!:!line!:!col!"

            echo   Found issue at !pos!
            echo   Applying fix...

            REM Apply code action
            gopls codeaction -w -exec "!pos!"

            if !errorlevel! equ 0 (
                echo   Fixed successfully
            ) else (
                echo   Warning: Fix may have failed
            )
            echo.
        )
    )
)

echo Done!

@echo off
setlocal enabledelayedexpansion

echo Finding and fixing unnecessary type arguments...
echo.

REM Find all Go files recursively
for /r %%f in (*.go) do (
    echo Checking %%f...
    
    REM Run gopls check and capture output
    for /f "tokens=1,2,3,4 delims=:" %%a in ('gopls check --severity=hint "%%f" 2^>^&1 ^| findstr /C:"unnecessary type arguments" ^| sort /R') do (
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

echo Done!
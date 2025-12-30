@echo off
cd /d "%~dp0frontend"
echo Starting frontend development server...
echo Working directory: %cd%
call npm run dev
pause

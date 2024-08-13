@echo off
set APP_ENV=dev
nodemon --watch "./**/*.go" --exec "go run main.go" --ext go
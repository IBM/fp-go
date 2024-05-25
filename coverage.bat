@echo off
go tool cover -html=build/cover.out -o build/cover.html
cov-report -ex ".*/cli/.*.go|.*/gen.go|.*/binds.go" build\cover.out

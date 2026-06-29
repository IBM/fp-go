@echo off

pushd ..\gen\v2
go generate .\...
popd

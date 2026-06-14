@echo off

pushd gen
go generate .\...
popd

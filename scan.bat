@echo off

busybox find . -type f -name "*\.go" | busybox xargs gopls check
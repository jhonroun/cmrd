SHELL := cmd.exe
.DEFAULT_GOAL := help

.PHONY: debug run build

debug:
	@echo Run with debug...
	go run cmd\app.go --debug

run:
	@echo Run...
	go run cmd\app.go

build:
	@echo Build and run...
	go build -o .\bin\cmrd.exe cmd\app.go

clean:
	@echo Clean...
	@rm -f .\bin\cmrd.exe
	@echo File 'cmrd.exe' deleted.

help:
	@echo Welcome to Cloud Mail.ru Downloader (CMRD)
	@echo release info: CMRD v0.1.0
	@echo
	@echo Copyright (c) 2024 LemTech (aka Jhon Roun)
	@echo repo: https://github.com/jhonroun/cmrd
	@echo
	@echo Usage:
	@echo   make debug        - run with debug
	@echo   make run          - just run app
	@echo   make build        - build and run app
	@echo   make clean        - clean all generated files

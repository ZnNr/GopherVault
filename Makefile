SHELL=bash
APP_VERSION=v1.0.0
.PHONY: install stop

install:
	docker-compose up --detach
	sleep 3
	go install -ldflags="-X 'github.com/ZnNr/GopherVault/cmd.version=$(APP_VERSION)' -X 'github.com/ZnNr/GopherVault/cmd.buildDate=$(shell date)'"

stop:
	docker-compose down
	docker image rm GopherVault-server --force & docker image rm GopherVault-migrate --force & docker image rm GopherVault-server --force
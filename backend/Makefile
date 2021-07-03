.PHONY: build
build:
	cd app && go build -o ../catalog-server main.go

.PHONY: run
run:
	cd app && go build -o ../catalog-server main.go && cd .. && ./catalog-server

.PHONY: migrate
migrate:
	cd app && go build -o ../migration migration.go && cd .. && ./migration

.DEFAULT_GOAL: build
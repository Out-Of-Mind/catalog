.PHONY: build
build:
	cd app && go build -o ../catalog-server main.go

.PHONY: run
run:
	cd app && go build -o ../catalog-server main.go && cd .. && ./catalog-server

.PHONY: migrate
migrate:
	cd app && go build -o ../migration migration.go && cd .. && ./migration

.PHONY: fmt
fmt:
	gofmt -l -s -w .

.PHONY: all
all:
	build
	migrate
	./catalog-server

.DEFAULT_GOAL: build
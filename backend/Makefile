.PHONY: build
build:
	cd app && go build -o ../catalog-server main.go

.PHONY: run
run:
	cd app && go build -o ../catalog-server main.go && cd .. && ./catalog-server

.DEFAULT_GOAL: build
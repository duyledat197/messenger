PROJECT_NAME := test
PKG := github.com/$(PROJECT_NAME)
MOD := $(GOPATH)/pkg/mod
COMPOSE_FILE := ./development/docker-compose.yml

compose:
	docker compose -f ${COMPOSE_FILE} up -d
gen-proto:
	docker compose -f ${COMPOSE_FILE} up generate_pb_go --build

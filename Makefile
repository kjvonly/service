SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

BASE_IMAGE_NAME := kjvonly/service
SERVICE_NAME    := bible-api
VERSION         := v0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ==============================================================================
# Hitting endpoints
es-search-local:
	curl -X POST  --data '{"query": "SELECT count(*) as matches from kjvonly where text like '\''%money%'\''"}' http://localhost:8080/v1/BibleSearchService.Search


# ==============================================================================
# Administration

migrate:
	go run tooling/kjvonly-admin/main.go migrate

seed:
	go run tooling/tooling//main.go seed


.PHONY: service

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.



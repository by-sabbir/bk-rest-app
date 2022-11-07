.PHONY:

DOCKER_REGISTRY ?= registry.sabbir.dev

get-deps:
	go mod tidy

test:
	docker compose up -d testdb && sleep 3 && \
	go test -coverprofile cover.out ./... -v -cover

show-coverage:
	go tool cover -html=cover.out

up:
	docker compose up -d --build api

logs:
	docker compose logs -f api

cleanup:
	docker compose down -v

login:
	docker login -u $$DOCKER_USER -p $$DOCKER_PASSWORD $(DOCKER_REGISTRY)

logout:
	docker logout $(DOCKER_REGISTRY)

report:
	go test -coverprofile=coverage.out ./internal/company/... && \
	gocover-cobertura < coverage.out > coverage.xml && \
	rm *.out
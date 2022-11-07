.PHONY:

DOCKER_REGISTRY ?= registry.sabbir.dev

get-deps:
	go mod tidy

test:
	go test -coverprofile=coverage.out ./internal/company/... -v -cover

show-coverage:
	go tool cover -html=coverage.out

publish:
	docker compose build api && \
	docker compose push api

logs:
	docker compose logs -f api

cleanup:
	docker compose down -v

login:
	docker login -u $$DOCKER_USER -p $$DOCKER_PASSWORD $(DOCKER_REGISTRY)

logout:
	docker logout $(DOCKER_REGISTRY)

report:
	gocover-cobertura < coverage.out > coverage.xml && \
	rm *.out
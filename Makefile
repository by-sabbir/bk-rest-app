.PHONY:

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

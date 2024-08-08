docker-compose-db:
	docker compose exec -it db psql -U root -d lenslocked

goose:
	goose -dir migrations postgres "host=localhost port=5432 user=root password=110963 dbname=lenslocked sslmode=disable" $(GOOSE_COMMAND)

GOOSE_COMMAND := $(if $(filter-out goose,$(MAKECMDGOALS)),$(filter-out goose,$(MAKECMDGOALS)),)

goose-create:
	goose -s -dir migrations create $(NAME) $(TYPE)

NAME := $(if $(filter-out goose-create,$(MAKECMDGOALS)),$(filter-out goose-create,$(MAKECMDGOALS)),)
TYPE := sql

%:
	@true

tailwind:
	npx tailwindcss -i tailwind/styles.css -o ../assets/styles.css --watch

docker-compose-up-and-build:
	docker compose -f docker-compose.yml -f docker-compose.production.yml up --build

.PHONY: docker-compose-db goose goose-create

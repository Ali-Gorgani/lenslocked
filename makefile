docker-compose-db:
	docker compose exec -it db psql -U root -d lenslocked

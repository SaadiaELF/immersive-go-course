.PHONY: up down build prune

prune:
	docker image prune -f

build:
	docker compose build

run: build
	docker compose up	

stop:
	docker compose down

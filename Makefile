.PHONY: wipe up

# Development
wipe:
	@echo 'Stopping and removing containers...'
	docker compose down
	docker rmi -f $(shell docker image ls -q ecom*)
up:
	docker compose up


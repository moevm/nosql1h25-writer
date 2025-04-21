.PHONY: compose-up
compose-up:
	docker compose up -d --build

.PHONY: compose-down
compose-down:
	docker compose down --remove-orphans

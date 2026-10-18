include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@printf "Clean all volume files? Risk of data loss. [y/N]: "; \
	read ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Environment cleaned"; \
	else \
		echo "Cleanup cancelled"; \
	fi


env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "A required parameter is missing seq. Example: make migrate-create seq=init"; \
		exit 1; \
	fi
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "A required parameter is missing action. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		-path=/migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"
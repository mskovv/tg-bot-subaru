up:
	@docker compose up -d

down:
	@docker compose down

cp-env:
	@cp .env.dist .env
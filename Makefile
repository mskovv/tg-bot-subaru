up:
	@docker compose up -d --build

down:
	@docker compose down

cp-env:
	@cp .env.dist .env

logs:
	@docker compose logs -f web

env:
	@docker compose exec -it web bash
start:
	@docker compose up --build

stop:
	@docker compose rm -v --force --stop
	@docker rm ticketeer
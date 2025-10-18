.PHONY: dev

dev:
	@trap 'kill 0' SIGINT SIGTERM; air & npm run dev & wait

.PHONY: reset

reset:
	@goose -dir migrations postgres "host=localhost port=5432 user=admin password=admin dbname=pixelparade sslmode=disable" reset

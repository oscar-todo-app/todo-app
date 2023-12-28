build:
	export TODO_DB_DSN=postgres://postgres:pa55word@127.0.0.1:5432/tasks?sslmode=disable
	sudo docker compose build
	sudo docker compose up -d
test:
	go test ./cmd

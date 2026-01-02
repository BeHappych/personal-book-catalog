start: run-docker run-go ## Полный запуск системы

run-docker: ## Запустить БД в Docker
	docker-compose up -d postgres

run-go: ## Запустить Go приложение
	go run ./cmd/main.go

stop: ## Остановить Docker
	docker-compose down
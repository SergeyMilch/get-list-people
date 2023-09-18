Название проекта:

"FIO Enrichment Service"

Краткое описание проекта:

"FIO Enrichment Service" - это сервис, который получает поток данных с ФИО через очередь Kafka, обогащает информацию о возрасте, поле и национальности, сохраняет данные в базу данных PostgreSQL и предоставляет REST и GraphQL интерфейсы для запросов, а также поддерживает кэширование в Redis.

Функциональность:

1.Слушает очередь Kafka с данными ФИО.
2.Обогащает данные возрастом, полом и национальностью.
3.Сохраняет данные в базу данных PostgreSQL.
4.Предоставляет REST API и GraphQL API для выполнения операций над данными.
5.Реализует кэширование данных с использованием Redis.

Запуск проекта

Убедитесь, что у вас установлены следующие зависимости:

1.PostgreSQL
2.Redis
3.Kafka
4.Go
5.Установите все необходимые библиотеки с помощью команды:
- go mod download

Создайте файл .env в корневой директории проекта и заполните его значениями:
POSTGRES_URL=postgres://user:password@localhost:5432/your_database?sslmode=disable&search_path=public
REDIS_URL=localhost:6379
KAFKA_BROKER_URL=localhost:9092
KAFKA_TOPIC=NAME_TOPIC
KAFKA_FAILED=NAME_TOPIC_FAILED
REDIS_PASSWORD=password
REDIS_DB=0

Запустите программу:
- go run main.go

Слои приложения 

router : описание маршрутов
handlers : работа с http 
services : бизнес-логика 
repositories : работа с бд через gorm 
models : модели базы данных 

ТРЕБОВАНИЯ ДЛЯ ЗАПУСКА 
Наличие Docker, Postgres запускается в контейнере, миграции применяются при запуске 

# ЗАПУСК ПРОЕКТА 
после клонирования в корне проекта: docker compose up --build 

Для проверки работы: localhost:{yourport(8080 by default)}/health 
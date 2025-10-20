# User API - REST API на Go с Gin и PostgreSQL

REST API для управления пользователями с поддержкой CRUD операций, пагинацией, фильтрацией и валидацией данных.

## Особенности

- Полный CRUD для пользователей
- PostgreSQL с использованием sqlx
- Docker и Docker Compose для развертывания
- Валидация входных данных
- Пагинация и фильтрация результатов
- Централизованная обработка ошибок
- Unit тесты
- Чистая архитектура (Repository → Service → Handler)
- Веб-интерфейс для управления пользователями

## Технологии

- Go 1.25.1
- Gin Web Framework
- PostgreSQL 16
- sqlx для работы с базой данных
- Docker & Docker Compose
- validator/v10 для валидации
- Go embed для встраивания статических файлов

## Структура проекта

```
user-api/
├── cmd/
│   └── api/
│       ├── main.go          # Точка входа приложения
│       └── static/          # Статические файлы (HTML)
│           └── index.html
├── internal/
│   ├── models/              # Модели данных
│   ├── repository/          # Слой работы с БД
│   ├── service/             # Бизнес-логика
│   ├── handlers/            # HTTP обработчики
│   ├── middleware/          # Middleware
│   └── database/            # Настройка подключения к БД
├── tests/                   # Тесты
├── migrations/              # SQL миграции
│   └── 001_create_users_table.sql
├── docker-compose.yml
├── Dockerfile
├── .env
├── go.mod
└── go.sum
```

## Установка и запуск

### Запуск через Docker Compose

```bash
# Клонировать репозиторий
git clone <your-repo>
cd user-api

# Запустить все сервисы
docker-compose up --build
```

API будет доступен по адресу `http://localhost:8080`

### Запуск локально

```bash
# Установить зависимости
go mod download

# Запустить PostgreSQL через Docker
docker-compose up postgres

# Запустить API
go run cmd/api/main.go
```

## API Endpoints

### Получить список пользователей

```bash
GET /api/v1/users?page=1&page_size=10&name=John&min_age=18&max_age=65
```

**Параметры запроса:**
- `page` - номер страницы (по умолчанию 1)
- `page_size` - размер страницы (по умолчанию 10, максимум 100)
- `name` - фильтр по имени
- `email` - фильтр по email
- `min_age` - минимальный возраст
- `max_age` - максимальный возраст

**Ответ:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 30,
      "created_at": "2025-10-20T18:00:00Z",
      "updated_at": "2025-10-20T18:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### Получить пользователя по ID

```bash
GET /api/v1/users/{id}
```

### Создать пользователя

```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

### Обновить пользователя

```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.new@example.com",
  "age": 31
}
```

### Удалить пользователя

```bash
DELETE /api/v1/users/{id}
```

**Ответ:** 204 No Content

### Health Check

```bash
GET /health
```

**Ответ:**
```json
{
  "status": "ok"
}
```

### Веб-интерфейс

```bash
GET /
```

Открывает веб-интерфейс для управления пользователями через браузер.

## Примеры использования

```bash
# Создать пользователя
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'

# Получить всех пользователей
curl http://localhost:8080/api/v1/users

# Получить пользователя по ID
curl http://localhost:8080/api/v1/users/1

# Обновить пользователя
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","age":26}'

# Удалить пользователя
curl -X DELETE http://localhost:8080/api/v1/users/1

# Фильтрация по параметрам
curl "http://localhost:8080/api/v1/users?page=1&page_size=5&name=Alice&min_age=20&max_age=30"

# Проверка состояния
curl http://localhost:8080/health
```

## Валидация данных

Автоматическая валидация при создании и обновлении пользователей:

**name:**
- Обязательное поле
- От 2 до 100 символов

**email:**
- Обязательное поле
- Валидный формат email

**age:**
- Обязательное поле
- От 1 до 150

**Пример ошибки валидации:**
```json
{
  "error": "Validation error",
  "message": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

## Фильтрация

Доступные параметры фильтрации:

- `name` - поиск по имени (регистронезависимый, частичное совпадение)
- `email` - поиск по email (регистронезависимый, частичное совпадение)
- `min_age` - минимальный возраст (включительно)
- `max_age` - максимальный возраст (включительно)

**Примеры:**
```bash
# Найти пользователей с именем содержащим "John"
curl "http://localhost:8080/api/v1/users?name=John"

# Найти пользователей от 20 до 30 лет
curl "http://localhost:8080/api/v1/users?min_age=20&max_age=30"

# Комбинированный фильтр
curl "http://localhost:8080/api/v1/users?name=Alice&min_age=25"
```

## Пагинация

Параметры:
- `page` - номер страницы (по умолчанию 1)
- `page_size` - элементов на странице (по умолчанию 10, максимум 100)

Ответ включает метаданные:
- `users` - массив пользователей
- `total` - общее количество
- `page` - текущая страница
- `page_size` - размер страницы
- `total_pages` - всего страниц

## Обработка ошибок

API возвращает структурированные ошибки в формате JSON.

**HTTP коды:**
- `200 OK` - успешный GET/PUT запрос
- `201 Created` - пользователь создан
- `204 No Content` - пользователь удален
- `400 Bad Request` - ошибка валидации
- `404 Not Found` - пользователь не найден
- `500 Internal Server Error` - ошибка сервера

## Docker

### Команды управления

```bash
# Запустить все сервисы
docker-compose up --build

# Запустить в фоновом режиме
docker-compose up -d

# Остановить сервисы
docker-compose down

# Остановить и удалить данные БД
docker-compose down -v

# Просмотр логов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f api
```

### Подключение к PostgreSQL

```bash
# Через docker-compose
docker-compose exec postgres psql -U postgres -d userdb

# Напрямую
docker exec -it user_api_db psql -U postgres -d userdb
```

**Полезные SQL команды:**
```sql
-- Показать все таблицы
\dt

-- Структура таблицы users
\d users

-- Все пользователи
SELECT * FROM users;

-- Очистить таблицу
TRUNCATE users RESTART IDENTITY CASCADE;
```

## Переменные окружения

Файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=userdb
SERVER_PORT=8080
```

## Тестирование

```bash
# Запустить все тесты
go test ./tests/...

# С покрытием
go test ./tests/... -cover

# С подробным выводом
go test ./tests/... -v

# Конкретный тест
go test ./tests/... -run TestCreateUser
```

## Миграции базы данных

SQL миграции находятся в папке `migrations/` и автоматически применяются при запуске PostgreSQL контейнера через `docker-entrypoint-initdb.d`.

Миграция `001_create_users_table.sql`:
- Создает таблицу users
- Добавляет индексы для оптимизации
- Устанавливает ограничения

## Архитектура

Проект следует принципам чистой архитектуры:

1. **Models** - структуры данных и DTO
2. **Repository** - работа с базой данных
3. **Service** - бизнес-логика
4. **Handlers** - HTTP обработчики
5. **Middleware** - промежуточные обработчики

Преимущества:
- Разделение ответственности
- Легкое тестирование
- Возможность замены компонентов
- Масштабируемость

## Зависимости

Основные библиотеки:

```go
github.com/gin-gonic/gin              // Web framework
github.com/lib/pq                     // PostgreSQL driver
github.com/jmoiron/sqlx               // SQL extensions
github.com/go-playground/validator    // Validation
github.com/joho/godotenv             // .env file support
github.com/stretchr/testify          // Testing toolkit
```

## Разработка

### Добавление нового поля в модель User

1. Обновить структуру в `internal/models/user.go`
2. Создать миграцию в `migrations/`
3. Обновить методы repository
4. Обновить валидацию при необходимости
5. Добавить тесты

### Добавление нового endpoint

1. Добавить метод в интерфейс `UserService`
2. Реализовать метод в `userService`
3. Добавить handler в `UserHandler`
4. Зарегистрировать route в `main.go`
5. Написать тесты

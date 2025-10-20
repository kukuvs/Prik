# User API - REST API на Go с Gin и PostgreSQL

REST API для управления пользователями с полной поддержкой CRUD операций, пагинацией, фильтрацией и валидацией данных.[1][2]

## Особенности

- ✅ Полный CRUD для пользователей
- ✅ PostgreSQL с использованием sqlx
- ✅ Docker и Docker Compose
- ✅ Валидация данных
- ✅ Пагинация и фильтрация
- ✅ Централизованная обработка ошибок
- ✅ Unit тесты
- ✅ Чистая архитектура (Repository → Service → Handler)

## Технологии

- Go 1.22
- Gin Web Framework
- PostgreSQL 16
- sqlx для работы с БД
- Docker & Docker Compose
- validator/v10 для валидации

## Структура проекта

```
user-api/
├── cmd/api/            # Точка входа
├── internal/
│   ├── models/         # Модели данных
│   ├── repository/     # Слой работы с БД
│   ├── service/        # Бизнес-логика
│   ├── handlers/       # HTTP обработчики
│   ├── middleware/     # Middleware
│   └── database/       # Настройка БД
├── tests/              # Тесты
├── migrations/         # SQL миграции
└── docker-compose.yml
```

## Установка и запуск

### 1. Клонирование проекта

```bash
git clone <your-repo>
cd user-api
```

### 2. Запуск через Docker Compose

```bash
docker-compose up --build
```

API будет доступен на `http://localhost:8080`

### 3. Запуск локально

```bash
# Установка зависимостей
go mod download

# Запуск PostgreSQL
docker-compose up postgres

# Запуск API
go run cmd/api/main.go
```

## API Endpoints

### Получить всех пользователей
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

**Ответ:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "created_at": "2025-10-20T18:00:00Z",
  "updated_at": "2025-10-20T18:00:00Z"
}
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

**Ответ (201 Created):**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "created_at": "2025-10-20T18:00:00Z",
  "updated_at": "2025-10-20T18:00:00Z"
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

**Ответ (200 OK):**
```json
{
  "id": 1,
  "name": "John Updated",
  "email": "john.new@example.com",
  "age": 31,
  "created_at": "2025-10-20T18:00:00Z",
  "updated_at": "2025-10-20T18:55:00Z"
}
```

### Удалить пользователя
```bash
DELETE /api/v1/users/{id}
```

**Ответ (204 No Content)**

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

## Примеры запросов с curl

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

# Фильтрация и пагинация
curl "http://localhost:8080/api/v1/users?page=1&page_size=5&name=Alice&min_age=20&max_age=30"

# Health check
curl http://localhost:8080/health
```

## Тестирование

```bash
# Запуск всех тестов
go test ./tests/...

# Запуск с покрытием
go test ./tests/... -cover

# Запуск с подробным выводом
go test ./tests/... -v

# Запуск конкретного теста
go test ./tests/... -run TestCreateUser
```

## Переменные окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=userdb
SERVER_PORT=8080
```

## Валидация данных

Автоматическая валидация при создании и обновлении пользователей:

- **name**: 
  - Обязательное поле
  - Минимум 2 символа
  - Максимум 100 символов

- **email**: 
  - Обязательное поле
  - Валидный формат email

- **age**: 
  - Обязательное поле
  - Минимум 1
  - Максимум 150

**Пример ошибки валидации:**
```json
{
  "error": "Validation error",
  "message": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

## Фильтрация

Поддерживаемые параметры фильтрации в GET /api/v1/users:

- `name` - поиск по имени (регистронезависимый, частичное совпадение)
- `email` - поиск по email (регистронезависимый, частичное совпадение)
- `min_age` - минимальный возраст (включительно)
- `max_age` - максимальный возраст (включительно)

**Примеры:**
```bash
# Найти всех пользователей с именем содержащим "John"
curl "http://localhost:8080/api/v1/users?name=John"

# Найти пользователей в возрасте от 20 до 30 лет
curl "http://localhost:8080/api/v1/users?min_age=20&max_age=30"

# Комбинированный фильтр
curl "http://localhost:8080/api/v1/users?name=Alice&min_age=25"
```

## Пагинация

Параметры пагинации:
- `page` - номер страницы (по умолчанию 1, минимум 1)
- `page_size` - количество элементов на странице (по умолчанию 10, максимум 100)

Ответ включает метаданные:
- `users` - массив пользователей на текущей странице
- `total` - общее количество пользователей
- `page` - текущая страница
- `page_size` - размер страницы
- `total_pages` - общее количество страниц

**Пример:**
```bash
# Получить вторую страницу по 20 элементов
curl "http://localhost:8080/api/v1/users?page=2&page_size=20"
```

## Обработка ошибок

API возвращает структурированные ошибки:

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

**HTTP коды ответов:**
- `200 OK` - успешный GET/PUT запрос
- `201 Created` - успешное создание пользователя
- `204 No Content` - успешное удаление
- `400 Bad Request` - ошибка валидации или некорректный запрос
- `404 Not Found` - пользователь не найден
- `500 Internal Server Error` - внутренняя ошибка сервера

## Docker

### Сборка и запуск

```bash
# Сборка и запуск всех сервисов
docker-compose up --build

# Запуск в фоновом режиме
docker-compose up -d

# Остановка сервисов
docker-compose down

# Остановка с удалением volumes (БД будет очищена)
docker-compose down -v

# Просмотр логов
docker-compose logs -f

# Просмотр логов конкретного сервиса
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

-- Показать структуру таблицы users
\d users

-- Показать всех пользователей
SELECT * FROM users;

-- Очистить таблицу
TRUNCATE users RESTART IDENTITY CASCADE;
```

## Миграции базы данных

SQL миграции находятся в папке `migrations/` и автоматически применяются при запуске PostgreSQL контейнера.

**migrations/001_create_users_table.sql:**
- Создает таблицу users
- Добавляет индексы для оптимизации поиска
- Устанавливает ограничения (constraints)

## Архитектура проекта

Проект следует принципам чистой архитектуры (Clean Architecture):

1. **Models** (`internal/models/`) - структуры данных и DTO
2. **Repository** (`internal/repository/`) - работа с базой данных
3. **Service** (`internal/service/`) - бизнес-логика
4. **Handlers** (`internal/handlers/`) - HTTP обработчики
5. **Middleware** (`internal/middleware/`) - промежуточные обработчики

**Преимущества:**
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

1. Обновите структуру в `internal/models/user.go`
2. Создайте новую миграцию в `migrations/`
3. Обновите repository методы
4. Обновите валидацию если необходимо
5. Обновите тесты

### Добавление нового endpoint

1. Добавьте метод в `UserService` интерфейс
2. Реализуйте метод в `userService`
3. Добавьте handler в `UserHandler`
4. Зарегистрируйте route в `main.go`
5. Добавьте тесты

## Лицензия

MIT License

## Автор

Лабораторная работа 8 - REST API на Go[6][1]

***

**Дата создания:** 20 октября 2025  
**Версия:** 1.0

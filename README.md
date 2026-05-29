# Student Distribution API

Backend-сервис на Go для распределения студентов по учебным группам.

---

## Что умеет сервис

- добавлять студентов через API;
- получать список студентов;
- удалять студента по ISU;
- импортировать студентов из CSV;
- создавать группы;
- получать список групп;
- запускать распределение студентов по группам;
- смотреть результат распределения в JSON;
- экспортировать результат распределения в CSV;
- Swagger UI по адресу `/swagger/`.

---

## Стек

- Go 1.26
- net/http (стандартная библиотека)
- PostgreSQL
- go-pg/pg v10
- swaggo/http-swagger — Swagger UI
- Docker Compose
- Python — генерация тестовых CSV-данных

---

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── docs/                     # сгенерированные swagger-файлы
├── internal/
│   ├── api/                  # слой конвертации моделей в API-ответы
│   ├── config/               # загрузка конфига из переменных окружения
│   ├── db/                   # подключение к PostgreSQL
│   ├── distributition/       # алгоритм распределения
│   ├── handler/              # HTTP-обработчики
│   ├── models/               # структуры данных (Student, Group, GroupDistribution)
│   ├── response/             # вспомогательная функция для JSON-ответов
│   └── storage/              # запросы к базе данных
├── migrations/               # SQL-файлы для создания таблиц
├── raw_data/                 # CSV-файлы для импорта/экспорта
├── scripts/                  # Python-скрипты для генерации тестовых данных
├── docker-compose.yml
├── Makefile
├── Makefile.mk.dist          # шаблон локальных переменных make
├── .env.dist                 # шаблон переменных окружения
└── go.mod
```

---

## Модель данных

### students

```sql
CREATE TABLE students (
    isu INTEGER PRIMARY KEY,
    full_name TEXT NOT NULL,
    telegram TEXT NOT NULL UNIQUE,
    score INTEGER NOT NULL CHECK (score >= 0),
    group_number INTEGER REFERENCES groups(number)
);
```

### groups

```sql
CREATE TABLE groups (
    number INTEGER PRIMARY KEY,
    capacity INTEGER NOT NULL CHECK (capacity > 0)
);
```

---

## Как работает распределение

1. Берём всех студентов из базы.
2. Сортируем по `score` по убыванию.
3. Берём все группы, сортируем по номеру.
4. Заполняем группы по порядку, пока не закончится вместимость или студенты.
5. Записываем каждому студенту `group_number`.

---

## Запуск

### 1. Инициализировать конфиг

```bash
make init
```

Создаёт `.env` и `Makefile.mk` из шаблонов. Заполни в них учётные данные PostgreSQL.

### 2. Поднять PostgreSQL

```bash
make db-up
```

### 3. Применить миграции

```bash
make migrate
```

Последовательно создаёт таблицы `students`, `groups` и добавляет внешний ключ.

### 4. Запустить API

```bash
make run
```

Сервер будет доступен на:

```text
http://localhost:8080
```

Проверка:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{"status":"ok"}
```

Swagger UI: `http://localhost:8080/swagger/`

---

## Переменные окружения

| Переменная          | По умолчанию | Описание                          |
|---------------------|--------------|-----------------------------------|
| `APP_PORT`          | `8080`       | Порт API                          |
| `POSTGRES_HOST`     | `localhost`  | Хост PostgreSQL                   |
| `POSTGRES_PORT`     | `5432`       | Порт PostgreSQL                   |
| `POSTGRES_USER`     | —            | Пользователь БД (обязательно)     |
| `POSTGRES_PASSWORD` | —            | Пароль БД (обязательно)           |
| `POSTGRES_DB`       | —            | Имя базы данных (обязательно)     |

Переменные читаются из окружения. При запуске через `make` они подхватываются из `Makefile.mk`.

---

## API

### Health check

```http
GET /health
```

---

## Students

### Создать студента

```http
POST /students
```

```bash
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{"isu":100001,"full_name":"Ivanov Ivan Ivanovich","telegram":"@student_100001","score":95}'
```

---

### Получить всех студентов

```http
GET /students
```

```bash
curl http://localhost:8080/students
```

---

### Удалить студента

```http
DELETE /students?isu=100001
```

```bash
curl -X DELETE "http://localhost:8080/students?isu=100001"
```

---

### Импортировать студентов из CSV

```http
POST /students/import
```

Формат CSV:

```csv
isu,full_name,telegram,score
100001,Ivanov Ivan Ivanovich,@student_100001,95
100002,Petrov Petr Petrovich,@student_100002,88
```

```bash
curl -X POST http://localhost:8080/students/import \
  -F "file=@raw_data/students.csv"
```

---

## Groups

### Создать группу

```http
POST /groups
```

```bash
curl -X POST http://localhost:8080/groups \
  -H "Content-Type: application/json" \
  -d '{"number":1,"capacity":25}'
```

---

### Получить все группы

```http
GET /groups
```

```bash
curl http://localhost:8080/groups
```

---

## Distribution

### Запустить распределение

```http
POST /distribution/run
```

```bash
curl -X POST http://localhost:8080/distribution/run
```

Возвращает список студентов с проставленными `group_number`.

---

### Получить результат распределения

```http
GET /distribution
```

```bash
curl http://localhost:8080/distribution
```

Пример ответа:

```json
[
  {
    "group_number": 1,
    "students": [
      {
        "isu": 100001,
        "full_name": "Ivanov Ivan Ivanovich",
        "telegram": "@student_100001",
        "score": 95,
        "group_number": 1
      }
    ]
  }
]
```

---

### Экспортировать распределение в CSV

```http
GET /distribution/export
```

```bash
curl http://localhost:8080/distribution/export \
  -o raw_data/distribution_export.csv
```

CSV-файл содержит колонки: `group_number`, `isu`, `full_name`, `telegram`, `score`.

---

## Генерация тестовых данных

```bash
make generate-csv   # сгенерировать CSV со студентами
make import-csv     # загрузить его в API
```

Файл сохраняется в `raw_data/students.csv`.

---

## Пример полного пайплайна

```bash
# Первый запуск
make init
# — заполни .env и Makefile.mk —
make db-up
make migrate

# Запустить API (первый терминал)
make run

# Во втором терминале
make generate-csv
make import-csv
make run-distribution
make export-distribution
```

Результат распределения: `raw_data/distribution_export.csv`

---

## Команды make

| Команда               | Описание                                    |
|-----------------------|---------------------------------------------|
| `make init`           | Создать `.env` и `Makefile.mk` из шаблонов  |
| `make db-up`          | Поднять PostgreSQL через Docker Compose     |
| `make db-down`        | Остановить PostgreSQL                       |
| `make migrate`        | Применить все миграции                      |
| `make run`            | Запустить API                               |
| `make generate-csv`   | Сгенерировать тестовый CSV                  |
| `make import-csv`     | Загрузить студентов из CSV                  |
| `make run-distribution` | Запустить распределение                   |
| `make export-distribution` | Сохранить результат в CSV            |
| `make swagger`        | Перегенерировать Swagger-документацию       |

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
- экспортировать результат распределения в CSV.

---

## Стек

- Go
- net/http
- PostgreSQL
- go-pg
- Docker Compose
- Python для генерации тестовых CSV-данных

---

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── db/
│   ├── handler/
│   ├── models/
│   └── response/
├── migrations/
├── raw_data/
├── scripts/
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

По папкам:

- `cmd/api` — точка входа приложения;
- `internal/db` — подключение к PostgreSQL;
- `internal/handler` — HTTP-обработчики;
- `internal/models` — структуры данных;
- `internal/response` — вспомогательная функция для JSON-ответов;
- `migrations` — SQL-файлы для создания таблиц;
- `scripts` — Python-скрипты для генерации тестовых данных;
- `raw_data` — CSV-файлы для импорта/экспорта.

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

Алгоритм сейчас простой:

1. Берём всех студентов из базы.
2. Сортируем их по `score` по убыванию.
3. Берём все группы и сортируем их по номеру.
4. Заполняем группы по порядку, пока не закончится вместимость или студенты.
5. Записываем каждому студенту `group_number`.
---

## Запуск

Сначала нужно поднять PostgreSQL:

```bash
make db-up
```

Потом запустить Go API:

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

Пример:

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

Загрузка:

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

Пример:

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

---

## Генерация тестовых данных

В проекте есть Python-скрипт, который генерирует большой CSV со студентами:

```bash
make generate-csv
```

Файл сохраняется сюда:

```text
raw_data/students.csv
```

После этого можно загрузить его в API:

```bash
make import-csv
```

---

## Пример полного пайплайна

В первом терминале:

```bash
make db-up
make run
```

Во втором терминале:

```bash
make generate-csv
make import-csv
make run-distribution
make export-distribution
```

После этого результат распределения будет лежать в:

```text
raw_data/distribution_export.csv
```

---

## Полезные команды

```bash
make db-up
```

Поднять PostgreSQL.

```bash
make db-down
```

Остановить PostgreSQL.

```bash
make db-reset
```

Полностью пересоздать базу.

```bash
make run
```

Запустить API.

```bash
make generate-csv
```

Сгенерировать тестовый CSV.

```bash
make import-csv
```

Загрузить студентов из CSV.

```bash
make run-distribution
```

Запустить распределение.

```bash
make export-distribution
```

Сохранить результат распределения в CSV.


---

```text
HTTP API → PostgreSQL → Docker → CSV import/export → простая бизнес-логика
```

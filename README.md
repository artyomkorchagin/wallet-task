# Сервис Кошельков
**Микросервис на Go для управления балансом кошельков с идемпотентностью, кешированием, PostgreSQL, Redis и корректной обработкой ошибок.**

---

## Требования
- Установлен Make
- Установлен Docker + Compose plugin


## Стек технологий

| Уровень         | Технология                     |
|----------------|--------------------------------|
| Язык           | Go 1.23+                       |
| Веб-фреймворк  | Gin     |
| База данных    | PostgreSQL       |
| Кеш            | Redis     |
| Логирование    | go.uber.org/zap             |
| Конфигурация   | spf13/vipe        |
| Тестирование   | testify, gomock, zaptest |
| Миграции       | goose |

---
## Makefile — Справочная таблица команд

| Команда       | Что делает                                                                | Зависит от             |
|---------------|---------------------------------------------------------------------------|------------------------|
| `make build`  | Собирает Docker-образы на основе `docker-compose.yml`                      | `config.env`           |
| `make up`     | Запускает все сервисы  в foreground                 | `config.env`           |
| `make down`   | Останавливает и удаляет контейнеры и сети (тома сохраняются)            | `config.env`           |
| `make restart`| Перезапускает всё: `down` `up`                                          | `down`, `up`           |
| `make test`   | Запускает все unit-тесты в проекте с подробным выводом               | —                      |
| `make cover`  | Генерирует HTML-отчёт о покрытии кода тестами `./tests/cover.html`              | —                      |
| `make clean`  | Полностью очищает окружение: контейнеры, тома, образы                        | `config.env`           |

---

## Быстрый старт

### 1. Клонируй и собери

```bash
git clone https://github.com/artyomkorchagin/wallet-task.git
cd wallet-task
make build
```
### 2. Запусти
```bash
make up
```

## Формат запросов
```bash
POST /api/v1/wallet
{
  "valletId": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "operationType": "DEPOSIT",
  "amount": 100,
}
GET /api/v1/wallet/:uuid
```
## Пример запросов
```bash
curl -X GET http://localhost:3000/api/v1/wallet/a1b2c3e4-5678-9012-3456-789012345678
curl -X POST http://localhost:3000/api/v1/wallet -H "Content-Type: application/json" -d "{\"valletId\": \"a1b2c3e4-5678-9012-3456-789012345678\", \"operationType\": \"WITHDRAW\", \"amount\": 100}"
```

## Заметки
- Чтобы обеспечить 1000 rps по одному кошельку нужно имплементировать очередь (планирую через редис)
- Нужно сделать документацию swagger
- ReferenceID должен быть сгенерирован на клиенте для обеспечение идемпотентности, но так как тут нет клиента, я генерирую ReferenceID внутри хендлера после получения запроса

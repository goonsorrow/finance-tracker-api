💰 Finance Tracker API
RESTful API сервис для трекинга личных финансов. Приложение позволяет пользователям регистрироваться, создавать кошельки и управлять транзакциями.

Написано на Go с использованием Clean Architecture.

✨ Features
JWT Authentication: Безопасная регистрация и вход (Access & Refresh токены).

Clean Architecture: Четкое разделение слоев (Handler → Service → Repository).

Graceful Shutdown: Корректное завершение работы сервера и БД при сигналах SIGINT/SIGTERM.

Context Propagation: Использование context.Context для отмены долгих запросов и таймаутов на всех уровнях.

Swagger UI: Автогенерируемая документация API.

Smart Configuration: Гибкое управление окружением через Makefile (Local vs Docker).

🛠 Технологический стек
Язык: Go (Golang)

Web Framework: Gin-gonic

Database: PostgreSQL (Driver: pgx / sqlx)

Config: Viper (поддержка .env и yaml)

Infrastructure: Docker & Docker Compose

Migrations: Golang-migrate

🚀 Быстрый старт
1. Клонирование и настройка
Склонируйте репозиторий и создайте файлы конфигурации из примеров одной командой:

bash
git clone https://github.com/goonsorrow/finance-tracker.git
cd finance-tracker
make init
(Команда make init создаст файлы .env.common, .env.docker, .env.local на основе шаблонов)

2. Запуск
Проект поддерживает два режима работы через Makefile:

Вариант А: Полный запуск в Docker (Рекомендуемый)
Поднимает приложение, базу данных и миграции в изолированных контейнерах.

bash
make docker
API будет доступно по адресу: http://localhost:8080

Вариант Б: Локальная разработка (Hybrid Mode)
Поднимает БД в Docker (на порту 5436, чтобы не конфликтовать с локальным Postgres), а само Go-приложение запускает локально.

bash
make run
Идеально для разработки, отладки и быстрой перекомпиляции.

📖 API Документация
После запуска сервера откройте Swagger UI:

👉 http://localhost:8080/swagger/index.html

Для авторизации используйте кнопку Authorize и введите токен в формате:
Bearer <ваш_access_token>

Основные эндпоинты:
```text
### Основные эндпоинты

#### 🔐 Auth (Авторизация)
*   `POST /auth/register` — Регистрация нового пользователя
*   `POST /auth/login` — Вход в систему (возвращает Access и Refresh токены)
*   `POST /auth/refresh` — Обновление Access токена

#### 💳 Wallets (Кошельки)
*   `POST /api/wallets/` — Создать новый кошелек
*   `GET /api/wallets/` — Список всех кошельков
*   `GET /api/wallets/{id}` — Детальная информация о кошельке
*   `PUT /api/wallets/{id}` — Редактировать кошелек
*   `DELETE /api/wallets/{id}` — Удалить кошелек

#### 💸 Records (Транзакции)
*   `POST /api/wallets/{id}/movements/` — Добавить транзакцию (доход/расход)
*   `GET /api/wallets/{id}/movements/` — История операций кошелька
*   `GET /api/wallets/{id}/movements/{trId}` — Детали транзакции
*   `PUT /api/wallets/{id}/movements/{trId}` — Изменить транзакцию
*   `DELETE /api/wallets/{id}/movements/{trId}` — Удалить транзакцию
```

📂 Структура проекта
```text
Проект следует принципам Clean Architecture:

├── cmd/                # Точка входа (main.go)
├── internal/
│   ├── handler/        # HTTP Handlers (Gin)
│   ├── service/        # Бизнес-логика
│   ├── repository/     # SQL запросы
│   └── models/         # Структуры данных
├── docs/               # Swagger файлы
├── schema/             # SQL миграции
├── .env.common         # Общие переменные (JWT, пароли)
├── .env.local          # Настройки для локального запуска
├── .env.docker         # Настройки для Docker сети
└── Makefile            # Команды автоматизации

```

## 🔜 Roadmap (Планы по развитию)
```text

### 1. 💱 Мультивалютность и конвертация
*   Интеграция с внешним API курсов валют (например, Fixer.io или OpenExchangeRates).
*   Эндпоинт `GET /auth/me` с отображением общего баланса пользователя, конвертированного в основную валюту.
*   Автоматическое обновление курсов в фоне (Worker pool).

### 2. ⚡ Оптимизация и Кэширование (Redis)
*   Внедрение **Redis** для кэширования частых GET-запросов (списки кошельков, транзакции).
*   Хранение Refresh-токенов в Redis для реализации мгновенного отзыва сессий (/logout).
*   Rate Limiting для защиты API от перегрузок.

### 3. 📊 Аналитика и Категории
*   Добавление сущности `Category` (Еда, Транспорт, Зарплата) с возможностью создавать свои категории.
*   Генерация отчетов по тратам за период (Monthly Spending Report).
*   SQL-агрегации для подсчета статистики по категориям.

```


🛠 Полезные команды Makefile

make test — Запуск тестов (будут добавлены)

make db-logs — Просмотр логов базы данных

make docker-down — Остановка и удаление контейнеров

make migrate-docker — Ручной накат миграций (делается автоматически при старте)

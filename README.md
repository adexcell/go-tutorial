```txt
.
├── cmd/
│   └── app/
│       └── main.go        # Точка входа
├── internal/              # Приватный код приложения
│   ├── config/            # Загрузка конфигов (Viper)
│   ├── handler/           # HTTP Handlers (Gin)
│   ├── service/           # Бизнес-логика
│   ├── repository/        # Работа с БД (Postgres)
│   └── domain/            # Интерфейсы и структуры данных
├── pkg/                   # Вспомогательный код (logger и т.д.)
├── migrations/            # SQL миграции
├── configs/               # YAML файлы
└── go.mod
```
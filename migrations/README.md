# Миграции базы данных

Этот каталог содержит все миграции базы данных для проекта HomeCloud DB Manager Service.

## Структура

```
migrations/
├── 000_create_database.sql      # Создание базы данных
├── 001_create_schema.sql        # Создание схемы
├── 002_create_users.sql         # Создание таблицы пользователей
├── 003_create_files_table.up.sql    # Создание таблицы файлов
├── 003_create_files_table.down.sql  # Откат таблицы файлов
├── 004_create_file_revisions_table.up.sql    # Создание таблицы версий файлов
├── 004_create_file_revisions_table.down.sql  # Откат таблицы версий файлов
├── 005_create_file_permissions_table.up.sql    # Создание таблицы прав доступа
├── 005_create_file_permissions_table.down.sql  # Откат таблицы прав доступа
├── migrate.sh                   # Скрипт управления миграциями
└── README.md                    # Эта документация
```

## Использование

### Через Makefile (рекомендуется)

```bash
# Применить все миграции
make migrate-up

# Показать статус миграций
make migrate-status

# Откатить последнюю миграцию
make migrate-down

# Откатить все миграции (с подтверждением)
make migrate-reset
```

### Напрямую через скрипт

```bash
# Применить все миграции
./migrations/migrate.sh up

# Показать статус
./migrations/migrate.sh status

# Откатить последнюю миграцию
./migrations/migrate.sh down

# Откатить все миграции
./migrations/migrate.sh reset

# Показать справку
./migrations/migrate.sh help
```

## Конфигурация

Скрипт использует следующие переменные окружения (с значениями по умолчанию):

- `DB_NAME=homecloud` - Имя базы данных
- `DB_USER=postgres` - Пользователь PostgreSQL
- `DB_HOST=localhost` - Хост базы данных
- `DB_PORT=5432` - Порт базы данных

Вы можете переопределить эти значения:

```bash
DB_HOST=my-db-server DB_USER=myuser ./migrations/migrate.sh up
```

## Отслеживание миграций

Скрипт автоматически создает таблицу `homecloud.migrations` для отслеживания примененных миграций. Это позволяет:

- Применять только новые миграции
- Откатывать миграции в правильном порядке
- Показывать статус каждой миграции

## Создание новых миграций

При создании новой миграции следуйте этим правилам:

1. Используйте последовательную нумерацию (006_, 007_, etc.)
2. Создавайте пару файлов: `.up.sql` и `.down.sql`
3. Файл `.up.sql` должен содержать изменения для применения
4. Файл `.down.sql` должен содержать изменения для отката

Пример:

```sql
-- 006_add_user_avatar.up.sql
ALTER TABLE homecloud.users ADD COLUMN avatar_url VARCHAR(255);

-- 006_add_user_avatar.down.sql
ALTER TABLE homecloud.users DROP COLUMN avatar_url;
```

## Безопасность

- Скрипт `migrate-reset` запрашивает подтверждение перед удалением всех данных
- Все операции выполняются в транзакциях
- При ошибке скрипт останавливается и не применяет дальнейшие миграции

## Требования

- PostgreSQL 12+
- Утилита `psql`
- Права на создание базы данных и таблиц 
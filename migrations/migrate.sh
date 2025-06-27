#!/bin/bash
set -e

CONFIG_FILE="config/config.local.yaml"

# Отладочная информация
echo "=== DEBUG INFO ==="
echo "CONFIG_FILE: $CONFIG_FILE"
echo "Current directory: $(pwd)"
echo "Files in current directory:"
ls -la
echo "Config file exists: $([ -f "$CONFIG_FILE" ] && echo "YES" || echo "NO")"
if [ -f "$CONFIG_FILE" ]; then
    echo "Config file contents:"
    cat "$CONFIG_FILE"
fi
echo "Environment variables:"
echo "DB_NAME: $DB_NAME"
echo "DB_USER: $DB_USER"
echo "DB_PASSWORD: $DB_PASSWORD"
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "=================="

# Функция для чтения значения из YAML
get_config_value() {
  local key="$1"
  grep "^\s*${key}:" "$CONFIG_FILE" | head -n1 | awk -F': ' '{print $2}' | tr -d '"'
}

# Если переменные окружения не заданы, читаем из config.local.yaml
DB_NAME="${DB_NAME:-$(get_config_value 'dbname')}"
DB_USER="${DB_USER:-$(get_config_value 'user')}"
DB_PASSWORD="${DB_PASSWORD:-$(get_config_value 'password')}"
DB_HOST="${DB_HOST:-$(get_config_value 'host')}"
DB_PORT="${DB_PORT:-$(get_config_value 'port')}"
MIGRATIONS_DIR="$(dirname "$0")"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции для вывода
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Функция для подключения к базе данных
psql_connect() {
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER "$@"
}

# Функция для создания таблицы миграций
create_migrations_table() {
    log_info "Создание таблицы для отслеживания миграций..."
    psql_connect $DB_NAME <<-EOF
        CREATE TABLE IF NOT EXISTS homecloud.migrations (
            id SERIAL PRIMARY KEY,
            filename VARCHAR(255) NOT NULL UNIQUE,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
EOF
    log_success "Таблица миграций создана"
}

# Функция для проверки, применена ли миграция
is_migration_applied() {
    local filename=$1
    local result=$(psql_connect $DB_NAME -t -c "SELECT COUNT(*) FROM homecloud.migrations WHERE filename = '$filename';")
    echo $result | tr -d ' '
}

# Функция для отметки миграции как примененной
mark_migration_applied() {
    local filename=$1
    psql_connect $DB_NAME -c "INSERT INTO homecloud.migrations (filename) VALUES ('$filename');"
}

# Функция для отметки миграции как отмененной
mark_migration_rolled_back() {
    local filename=$1
    psql_connect $DB_NAME -c "DELETE FROM homecloud.migrations WHERE filename = '$filename';"
}

# Функция для применения миграции
apply_migration() {
    local file=$1
    local filename=$(basename "$file")
    
    if [ "$(is_migration_applied "$filename")" = "0" ]; then
        log_info "Применение миграции: $filename"
        psql_connect $DB_NAME -f "$file"
        mark_migration_applied "$filename"
        log_success "Миграция $filename применена"
    else
        log_warning "Миграция $filename уже применена, пропускаем"
    fi
}

# Функция для отката миграции
rollback_migration() {
    local up_file=$1
    local down_file=${up_file/.up.sql/.down.sql}
    local filename=$(basename "$up_file")
    
    if [ "$(is_migration_applied "$filename")" = "1" ]; then
        if [ -f "$down_file" ]; then
            log_info "Откат миграции: $filename"
            psql_connect $DB_NAME -f "$down_file"
            mark_migration_rolled_back "$filename"
            log_success "Миграция $filename откачена"
        else
            log_error "Файл отката не найден: $down_file"
            exit 1
        fi
    else
        log_warning "Миграция $filename не применена, пропускаем"
    fi
}

# Функция для показа статуса миграций
show_status() {
    log_info "Статус миграций:"
    echo "----------------------------------------"
    
    for file in $(find "$MIGRATIONS_DIR" -name "*.up.sql" | sort); do
        local filename=$(basename "$file")
        local applied=$(is_migration_applied "$filename")
        
        if [ "$applied" = "1" ]; then
            echo -e "${GREEN}✓${NC} $filename"
        else
            echo -e "${RED}✗${NC} $filename"
        fi
    done
    echo "----------------------------------------"
}

# Функция для применения всех миграций
migrate_up() {
    log_info "Применение всех миграций..."
    
    # Проверка подключения к базе данных
    log_info "Проверка подключения к PostgreSQL..."
    if ! psql_connect -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
        log_error "Не удается подключиться к PostgreSQL"
        log_error "DB_HOST: $DB_HOST, DB_PORT: $DB_PORT, DB_USER: $DB_USER"
        exit 1
    fi
    log_success "Подключение к PostgreSQL успешно"
    
    # Создание базы данных и схемы
    log_info "Создание базы данных..."
    psql_connect -d postgres -f "$MIGRATIONS_DIR/000_create_database.up.sql"
    
    log_info "Создание схемы..."
    psql_connect -d "$DB_NAME" -f "$MIGRATIONS_DIR/001_create_schema.up.sql"
    
    # Создание таблицы миграций
    create_migrations_table
    
    # Применение всех .up.sql файлов
    for file in $(find "$MIGRATIONS_DIR" -name "*.up.sql" | sort); do
        # Пропускаем первые две миграции, так как они уже применены
        if [[ "$file" == *"000_create_database.up.sql" ]] || [[ "$file" == *"001_create_schema.up.sql" ]]; then
            continue
        fi
        apply_migration "$file"
    done
    
    log_success "Все миграции применены"
}

# Функция для отката последней миграции
migrate_down() {
    log_info "Откат последней миграции..."
    
    # Получаем последнюю примененную миграцию
    local last_migration=$(psql_connect $DB_NAME -t -c "SELECT filename FROM homecloud.migrations ORDER BY applied_at DESC LIMIT 1;" | tr -d ' ')
    
    if [ -z "$last_migration" ]; then
        log_warning "Нет примененных миграций для отката"
        return
    fi
    
    local up_file="$MIGRATIONS_DIR/$last_migration"
    rollback_migration "$up_file"
}

# Функция для отката всех миграций
migrate_reset() {
    log_warning "Откат всех миграций..."
    read -p "Вы уверены? Это удалит все данные! (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        for file in $(find "$MIGRATIONS_DIR" -name "*.up.sql" | sort -r); do
            local filename=$(basename "$file")
            if [ "$(is_migration_applied "$filename")" = "1" ]; then
                rollback_migration "$file"
            fi
        done
        log_success "Все миграции откачены"
    else
        log_info "Операция отменена"
    fi
}

# Основная логика
case "${1:-up}" in
    "up")
        migrate_up
        ;;
    "down")
        migrate_down
        ;;
    "reset")
        migrate_reset
        ;;
    "status")
        show_status
        ;;
    "help"|"-h"|"--help")
        echo "Использование: $0 [команда]"
        echo ""
        echo "Команды:"
        echo "  up     - Применить все миграции (по умолчанию)"
        echo "  down   - Откатить последнюю миграцию"
        echo "  reset  - Откатить все миграции"
        echo "  status - Показать статус миграций"
        echo "  help   - Показать эту справку"
        echo ""
        echo "Переменные окружения:"
        echo "  DB_NAME - Имя базы данных (по умолчанию: homecloud)"
        echo "  DB_USER - Пользователь БД (по умолчанию: postgres)"
        echo "  DB_HOST - Хост БД (по умолчанию: localhost)"
        echo "  DB_PORT - Порт БД (по умолчанию: 5432)"
        ;;
    *)
        log_error "Неизвестная команда: $1"
        echo "Используйте '$0 help' для справки"
        exit 1
        ;;
esac 
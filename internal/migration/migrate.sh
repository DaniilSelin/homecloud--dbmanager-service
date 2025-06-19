#!/bin/bash
set -e
DB="homecloud"
USER="postgres"

# Проверка наличия dblink
psql -U $USER -d postgres -c "CREATE EXTENSION IF NOT EXISTS dblink;"

# Создание базы данных, если не существует
psql -U $USER -d postgres -f internal/migration/000_create_database.sql

# Создание схемы и таблиц
psql -U $USER -d $DB -f internal/migration/000_create_schema.sql
psql -U $USER -d $DB -f internal/migration/001_create_users.sql 
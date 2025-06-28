-- Удаление таблицы пользователей
DROP TABLE IF EXISTS homecloud.users CASCADE;
 
-- Удаление функции триггера
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE; 
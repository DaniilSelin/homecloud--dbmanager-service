-- Создание базы данных homecloud
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'homecloud') THEN
        CREATE DATABASE homecloud;
    END IF;
END
$$; 
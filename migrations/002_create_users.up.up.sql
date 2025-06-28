-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS homecloud.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    is_email_verified BOOLEAN DEFAULT false,
    role VARCHAR(20) DEFAULT 'user',
    storage_quota BIGINT DEFAULT 10737418240, -- 10GB в байтах
    used_space BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_users_username ON homecloud.users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON homecloud.users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON homecloud.users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_role ON homecloud.users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON homecloud.users(is_active);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON homecloud.users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 
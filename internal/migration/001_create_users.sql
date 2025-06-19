CREATE TABLE IF NOT EXISTS dbmanager.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Идентификация
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    -- Аутентификация и безопасность
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMP,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    -- Информация о хранилище
    storage_quota BIGINT NOT NULL DEFAULT 10737418240, -- 10 GiB
    used_space BIGINT NOT NULL DEFAULT 0,
    -- Роли и разрешения
    role TEXT NOT NULL DEFAULT 'user', -- user / admin / readonly / etc
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    -- Метаданные
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
); 
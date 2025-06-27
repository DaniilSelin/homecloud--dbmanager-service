-- Создание таблицы прав доступа к файлам
CREATE TABLE homecloud.file_permissions (
    id           UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id      UUID      NOT NULL REFERENCES homecloud.files(id) ON DELETE CASCADE,
    grantee_id   UUID,     -- user|group|domain или NULL для «всех»
    grantee_type TEXT      NOT NULL,  -- USER, GROUP, DOMAIN, ANYONE
    role         TEXT      NOT NULL,  -- OWNER, ORGANIZER, FILE_OWNER, WRITER, COMMENTER, READER
    allow_share  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP NOT NULL DEFAULT now()
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_file_permissions_file_id ON homecloud.file_permissions(file_id);
CREATE INDEX idx_file_permissions_grantee_id ON homecloud.file_permissions(grantee_id);
CREATE INDEX idx_file_permissions_grantee_type ON homecloud.file_permissions(grantee_type);
CREATE INDEX idx_file_permissions_role ON homecloud.file_permissions(role);
CREATE INDEX idx_file_permissions_created_at ON homecloud.file_permissions(created_at);

-- Уникальный индекс для предотвращения дублирования прав
CREATE UNIQUE INDEX idx_file_permissions_unique ON homecloud.file_permissions(file_id, grantee_id, grantee_type);

-- Ограничения для валидации данных
ALTER TABLE homecloud.file_permissions ADD CONSTRAINT chk_grantee_type 
    CHECK (grantee_type IN ('USER', 'GROUP', 'DOMAIN', 'ANYONE'));

ALTER TABLE homecloud.file_permissions ADD CONSTRAINT chk_role 
    CHECK (role IN ('OWNER', 'ORGANIZER', 'FILE_OWNER', 'WRITER', 'COMMENTER', 'READER')); 
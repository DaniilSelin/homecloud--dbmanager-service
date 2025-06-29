-- Удаляем старый индекс
DROP INDEX IF EXISTS homecloud.idx_files_unique_name_root;

-- Создаем новый индекс, который учитывает owner_id и parent_id
CREATE UNIQUE INDEX idx_files_unique_name_owner_parent ON homecloud.files (name, owner_id, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000')); 
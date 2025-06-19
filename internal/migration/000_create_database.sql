DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'homecloud') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE homecloud');
   END IF;
END$$; 
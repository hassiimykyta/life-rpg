\echo '>>> INIT 001-auth.sql STARTED <<<'

-- 1) создаём роль, если нет
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'auth_app') THEN
    CREATE ROLE auth_app LOGIN PASSWORD 'auth';
  ELSE
    ALTER ROLE auth_app WITH LOGIN PASSWORD 'auth';
  END IF;
END$$;

-- 2) создаём БД, если нет (НЕ в DO; используем \gexec)
SELECT 'CREATE DATABASE auth_db OWNER auth_app'
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'auth_db')\gexec

-- 3) права на public-схему (после создания БД)
\connect auth_db
ALTER SCHEMA public OWNER TO auth_app;
GRANT ALL ON SCHEMA public TO auth_app;

\echo '>>> INIT 001-auth.sql FINISHED <<<'
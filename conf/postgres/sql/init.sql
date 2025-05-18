DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'pg') THEN
        CREATE DATABASE pg;
    END IF;
END
$$;

DO
$$
    BEGIN
        CREATE SCHEMA IF NOT EXISTS todo;
        CREATE ROLE {{.Username }} WITH LOGIN ENCRYPTED PASSWORD {{ printf "'%s'" .Password }};
        GRANT CONNECT ON DATABASE todo TO {{ .Username }};
        GRANT USAGE ON SCHEMA todo TO {{ .Username }};
        ALTER DEFAULT PRIVILEGES IN SCHEMA todo GRANT SELECT, INSERT, DELETE, UPDATE ON TABLES TO {{ .Username }};
        ALTER DEFAULT PRIVILEGES IN SCHEMA todo GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO {{ .Username }};
        ALTER ROLE {{ .Username }} SET SEARCH_PATH TO todo;
    EXCEPTION
        WHEN OTHERS THEN
            RAISE NOTICE 'failed to create schema and role';
    END
$$;
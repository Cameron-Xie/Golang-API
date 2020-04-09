CREATE TABLE IF NOT EXISTS todo.tasks
(
    id          uuid                     not null primary key,
    name        varchar(150)             not null,
    description text                     not null,
    created_at  timestamp with time zone not null default now(),
    updated_at  timestamp with time zone
);

CREATE OR REPLACE FUNCTION set_updated_at_on_changes() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$$
BEGIN
    NEW.updated_at = OLD.updated_at;

    IF OLD IS DISTINCT FROM NEW THEN
        NEW.updated_at = NOW();
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS generate_updated_at_timestamp_before_task_updated ON todo.tasks;
CREATE TRIGGER generate_updated_at_timestamp_before_task_updated
    BEFORE UPDATE
    ON todo.tasks
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_on_changes();

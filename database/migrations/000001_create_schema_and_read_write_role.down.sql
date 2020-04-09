DROP OWNED BY {{ .Username }};
DROP ROLE IF EXISTS {{ .Username }};
DROP SCHEMA IF EXISTS todo;

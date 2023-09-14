CREATE TABLE IF NOT EXISTS people (
    id BIGSERIAL PRIMARY KEY,
    user_name TEXT,
    surname TEXT,
    patronymic TEXT,
    age INT,
    gender TEXT,
    nationality TEXT
);
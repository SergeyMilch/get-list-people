CREATE TABLE IF NOT EXISTS public.people (
    id BIGSERIAL PRIMARY KEY,
    user_name TEXT,
    surname TEXT,
    patronymic TEXT,
    age INT,
    gender TEXT,
    nationality TEXT
);
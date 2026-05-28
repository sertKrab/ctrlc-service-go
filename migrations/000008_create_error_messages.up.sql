CREATE TABLE error_messages (
    id         SERIAL PRIMARY KEY,
    code       VARCHAR(100) NOT NULL UNIQUE,
    locale_th  TEXT NOT NULL,
    locale_en  TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

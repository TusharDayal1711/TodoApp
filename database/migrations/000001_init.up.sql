CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, -- creating unique primary key using serial //
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP -- Null = soft delete
);

--create unique index use_email_idx(email,username,id)where archived_at is null;

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY, --unique id for using session_id as tokens
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,-- creating unique primary key using postres serial, //
    user_id INTEGER NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP -- Null = soft delete
);

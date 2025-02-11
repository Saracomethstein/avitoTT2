-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(256),
    balance INTEGER DEFAULT 1000
);

-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавление тестовых данных
INSERT INTO users (username, password) VALUES
    ('admin', crypt('password', gen_salt('bf'))),
    ('user', crypt('pass', gen_salt('bf')))
ON CONFLICT (username) DO NOTHING;


CREATE INDEX idx_users_username ON users(username);

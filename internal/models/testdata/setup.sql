DROP TABLE IF EXISTS snippets;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;


CREATE TABLE users (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    avatar_url VARCHAR(255) NOT NULL DEFAULT '/static/img/avatars/penguin.png'
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE TABLE snippets (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    visibility_level INTEGER NOT NULL DEFAULT 0
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created, avatar_url) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00',
    '/static/img/avatars/penguin.png'
);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

create table sessions (
	token char(43) primary key,
	data BLOB not null,
	expiry timestamp(6) not null
);

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	user_id int not null,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

create index snippets_created_idx on snippets (created);

create index sessions_expiry_idx on sessions (expiry);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
	

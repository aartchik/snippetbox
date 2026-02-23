create table if not exists users (
	id int not null primary key AUTO_INCREMENT,
	name varchar(255) not null,
	email varchar(255) not null,
	hashed_password char(60) not null,
	created datetime not null );

create table if not exists snippets (
	id int primary key AUTO_INCREMENT,
	title varchar(100) not null,
	content text not null,
	created datetime not null,
	expires datetime not null,
	user_id int not null );

create table if not exists sessions (
	token char(43) primary key,
	data BLOB not null,
	expiry timestamp(6) not null );


alter table users add constraint users_email_unique unique (email);

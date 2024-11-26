create table users (
	id serial primary key,
	username varchar(255) not null unique,
	password_hash text not null);

alter table users add column role varchar(50) not null default 'user';

create table refresh_tokens (
	id serial primary key,
	user_id integer not null references users(id) on delete cascade,
	token text not null unique,
	expires_at timestamp not null);

create table gateways (
	id serial primary key,
	name varchar(255) not null,
	ip inet not null unique
);

create table bilhetes (
	id serial primary key,
	bid varchar(255) unique not null,
	lega varchar(255),
	legb varchar(255),
	foreign key (lega) references call_records(call_id) on delete cascade,
	foreign key (legb) references call_records(call_id) on delete cascade,
	created_at timestamp default now()
);

